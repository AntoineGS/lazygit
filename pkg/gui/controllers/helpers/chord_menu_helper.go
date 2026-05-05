package helpers

import (
	"fmt"
	"sync"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// scheduler abstracts time.AfterFunc so tests can drive the timer
// deterministically.
type scheduler interface {
	Schedule(d time.Duration, fn func()) (cancel func())
}

type realScheduler struct{}

func (realScheduler) Schedule(d time.Duration, fn func()) func() {
	t := time.AfterFunc(d, fn)
	return func() { t.Stop() }
}

// ChordMenuHelper subscribes to gocui's chord-state callback and opens a
// menu showing chord continuations. Once the menu is open, the helper
// owns the prefix state; gocui chord state is used only during the
// optional pre-open delay window.
type ChordMenuHelper struct {
	c          *HelperCommon
	viewHelper *ViewHelper
	scheduler  scheduler

	mu          sync.Mutex
	cancelTimer func()
	menuOpen    bool
	opening     bool
	// timerGen increments on every chord-state change. A scheduled
	// timer captures the generation at schedule time and bails if it
	// has advanced — this protects against the documented race where
	// time.AfterFunc's callback was already queued on the UI thread
	// when Stop() was called, so the callback fires *after* the chord
	// was cancelled or replaced.
	timerGen int
	// deferredPrefix holds a non-empty chord-state change that arrived
	// while we were mid-open. Replayed via OnChordStateChange once
	// opening completes; processing it inline would race with the
	// in-flight c.Menu call. Empty prefixes (cancellations) are dropped
	// at receive time — see OnChordStateChange — so deferredPrefix == nil
	// uniquely means "no deferral".
	deferredPrefix []gocui.Key
	// originContext is preserved across the menu Push/Pop so continuations
	// remain scoped to the originating view even though gocui's chord
	// state is cleared by Push.
	originContext types.Context

	openHookForTest    func([]gocui.Key)
	closeHookForTest   func()
	refreshHookForTest func([]gocui.Key)

	// titleFuncs maps "<ctxName>::<prefixLabel>" to a dynamic title
	// resolver. Populated by helpers via RegisterTitleFunc; consulted by
	// titleForPrefix before falling back to the static group Name.
	titleFuncs map[string]func() string
}

func NewChordMenuHelper(c *HelperCommon, viewHelper *ViewHelper) *ChordMenuHelper {
	return &ChordMenuHelper{c: c, viewHelper: viewHelper, scheduler: realScheduler{}}
}

// RegisterTitleFunc registers a dynamic title resolver for a chord prefix
// in a given context. The function is called each time the chord popup
// opens or refreshes; it must be cheap (no I/O) and may be invoked from
// any goroutine.
func (h *ChordMenuHelper) RegisterTitleFunc(ctxName, prefixLabel string, fn func() string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.titleFuncs == nil {
		h.titleFuncs = map[string]func() string{}
	}
	h.titleFuncs[ctxName+"::"+prefixLabel] = fn
}

func (h *ChordMenuHelper) titleForPrefix(
	prefix []gocui.Key,
	ctxName string,
	groups map[string]map[string]config.KeybindingGroupConfig,
) string {
	label := config.LabelForKeySequence(prefix)
	h.mu.Lock()
	ctxFn := h.titleFuncs[ctxName+"::"+label]
	globalFn := h.titleFuncs["global::"+label]
	h.mu.Unlock()
	// Empty return falls through to the static label so callbacks can
	// opt out of dynamic naming for some states without re-implementing
	// the static fallback.
	if ctxFn != nil {
		if title := ctxFn(); title != "" {
			return title
		}
	}
	if globalFn != nil {
		if title := globalFn(); title != "" {
			return title
		}
	}
	if g, ok := config.LookupGroup(groups, ctxName, label); ok {
		return g.Name
	}
	return fmt.Sprintf("Chord: %s …", label)
}

// buildMenuItems maps chord-continuation rows into MenuItems.
//
// If any row has a non-empty inline tooltip, every row gets a 2nd
// LabelColumns entry — possibly empty — so column widths align.
//
// Rows whose key label fails to round-trip through KeyFromLabel are
// logged and dropped. Including them would produce a MenuItem with the
// zero gocui.Key — an unbindable row with no error surfaced.
func (h *ChordMenuHelper) buildMenuItems(
	infos []bindingInfo,
	prefix []gocui.Key,
	extendFn func([]gocui.Key) error,
) []*types.MenuItem {
	anyExtra := lo.SomeBy(infos, func(i bindingInfo) bool { return i.tooltip != "" })

	return lo.FilterMap(infos, func(info bindingInfo, _ int) (*types.MenuItem, bool) {
		nextKey, ok := config.KeyFromLabel(info.key)
		if !ok {
			h.c.Log.Errorf("ChordMenu buildMenuItems: unparseable key label %q; skipping row", info.key)
			return nil, false
		}

		cols := []string{info.description}
		if anyExtra {
			cols = append(cols, info.tooltip)
		}

		thisKey := nextKey.Sequence()[0]

		item := &types.MenuItem{
			LabelColumns: cols,
			Key:          thisKey,
		}

		if info.isGroup {
			extended := append([]gocui.Key{}, prefix...)
			extended = append(extended, thisKey)
			item.OnPress = func() error { return extendFn(extended) }
			return item, true
		}

		// Mirror GetDisabledReason onto the MenuItem so disabled rows
		// surface the standard toast on press, matching the regular
		// dispatch path. AllowFurtherDispatching rows are filtered out
		// upstream by BuildChordContinuations.
		if info.binding.GetDisabledReason != nil {
			item.DisabledReason = info.binding.GetDisabledReason()
		}
		item.OnPress = info.binding.Handler
		item.Tooltip = info.binding.GetTooltip()
		return item, true
	})
}

// OnChordStateChange is invoked by gocui when the chord state changes.
// A negative ChordPopupDelayMs disables the popup entirely.
//
// While `opening` is set: empty callbacks are dropped (these are the
// re-entrant Push callbacks fired by c.Menu); non-empty callbacks are
// queued and replayed after openMenu completes — running them inline
// would race with the in-flight c.Menu call.
func (h *ChordMenuHelper) OnChordStateChange(prefix []gocui.Key) {
	h.mu.Lock()
	if h.opening {
		if len(prefix) == 0 {
			h.mu.Unlock()
			return
		}
		h.deferredPrefix = append([]gocui.Key(nil), prefix...)
		h.mu.Unlock()
		return
	}
	delayMs := h.c.UserConfig().ChordPopupDelayMs
	h.timerGen++
	gen := h.timerGen
	if h.cancelTimer != nil {
		h.cancelTimer()
		h.cancelTimer = nil
	}
	wasOpen := h.menuOpen
	h.mu.Unlock()

	if delayMs < 0 {
		if wasOpen {
			h.closeMenu()
		}
		return
	}
	if len(prefix) == 0 {
		if wasOpen {
			h.closeMenu()
		}
		return
	}

	prefixCopy := append([]gocui.Key(nil), prefix...)
	if wasOpen {
		h.refreshMenu(prefixCopy)
		return
	}
	if delayMs == 0 {
		h.openMenu(prefixCopy)
		return
	}

	h.mu.Lock()
	h.cancelTimer = h.scheduler.Schedule(time.Duration(delayMs)*time.Millisecond, func() {
		h.mu.Lock()
		if h.timerGen != gen {
			h.mu.Unlock()
			return
		}
		h.mu.Unlock()
		h.openMenu(prefixCopy)
	})
	h.mu.Unlock()
}

// openMenu opens the chord menu via c.Menu.
//
// On first open we capture the chord-start context into originContext
// so re-opens for group continuations keep resolving against the
// originating view, not the menu context that becomes "current" after
// Push.
func (h *ChordMenuHelper) openMenu(prefix []gocui.Key) {
	if h.openHookForTest != nil {
		h.mu.Lock()
		h.opening = true
		h.menuOpen = true
		h.mu.Unlock()
		h.openHookForTest(prefix)
		h.mu.Lock()
		h.opening = false
		h.mu.Unlock()
		h.replayDeferred()
		return
	}

	h.mu.Lock()
	if h.originContext == nil {
		viewName := h.c.GocuiGui().PendingChordView()
		if viewName != "" {
			if ctx, ok := h.viewHelper.ContextForView(viewName); ok {
				h.originContext = ctx
			}
		}
		if h.originContext == nil {
			h.originContext = h.c.Context().Current()
		}
	}
	h.mu.Unlock()

	groups := h.c.UserConfig().KeybindingGroups
	ctxName := h.contextNameForChord()
	infos := BuildChordContinuations(h.gatherBindings(), prefix, groups, ctxName)
	items := h.buildMenuItems(infos, prefix, h.extendPrefix)

	h.mu.Lock()
	h.opening = true
	h.mu.Unlock()

	err := h.c.Menu(types.CreateMenuOptions{
		Title: h.titleForPrefix(prefix, ctxName, groups),
		Items: items,
	})

	h.mu.Lock()
	h.opening = false
	if err == nil {
		h.menuOpen = true
	}
	h.mu.Unlock()

	h.replayDeferred()
}

// replayDeferred drains a non-empty chord-state change captured by
// OnChordStateChange while opening was in progress. Runs after
// opening = false, so the replayed call follows the normal
// open/refresh paths without racing with c.Menu.
func (h *ChordMenuHelper) replayDeferred() {
	h.mu.Lock()
	if h.deferredPrefix == nil {
		h.mu.Unlock()
		return
	}
	prefix := h.deferredPrefix
	h.deferredPrefix = nil
	h.mu.Unlock()
	h.OnChordStateChange(prefix)
}

// refreshMenu swaps the open menu's items and title in place — no Push
// or Pop. Only the menu view's keybindings are re-registered (the
// per-item chord-continuation keys are the only ones that change
// row-to-row); all other controllers' bindings stay in place.
func (h *ChordMenuHelper) refreshMenu(newPrefix []gocui.Key) {
	if h.refreshHookForTest != nil {
		h.refreshHookForTest(newPrefix)
		return
	}
	groups := h.c.UserConfig().KeybindingGroups
	ctxName := h.contextNameForChord()
	infos := BuildChordContinuations(h.gatherBindings(), newPrefix, groups, ctxName)
	items := h.buildMenuItems(infos, newPrefix, h.extendPrefix)

	h.c.Contexts().Menu.SetMenuItems(items, nil)
	h.c.Views().Menu.Title = h.titleForPrefix(newPrefix, ctxName, groups)
	if err := h.c.RefreshMenuKeybindings(); err != nil {
		h.c.Log.Errorf("ChordMenu refreshMenu RefreshMenuKeybindings: %v", err)
	}
	h.c.PostRefreshUpdate(h.c.Contexts().Menu)
}

func (h *ChordMenuHelper) closeMenu() {
	h.mu.Lock()
	if !h.menuOpen {
		h.mu.Unlock()
		return
	}
	h.menuOpen = false
	h.originContext = nil
	h.mu.Unlock()

	if h.closeHookForTest != nil {
		h.closeHookForTest()
		return
	}
	h.c.Context().Pop()
}

// extendPrefix runs after MenuContext.OnMenuPress has popped the
// current menu, so we push a fresh one rather than refreshing in place.
// originContext is preserved across the Pop+re-Push so continuations
// resolve against the originating view.
func (h *ChordMenuHelper) extendPrefix(newPrefix []gocui.Key) error {
	h.mu.Lock()
	h.menuOpen = false
	h.mu.Unlock()
	h.openMenu(newPrefix)
	return nil
}

func (h *ChordMenuHelper) NotifyMenuClosed() {
	h.mu.Lock()
	h.menuOpen = false
	h.originContext = nil
	h.mu.Unlock()
}

func (h *ChordMenuHelper) IsOpen() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.menuOpen
}

// gatherBindings collects bindings from the chord-start view's context
// plus the global context, deduplicating by Key. We have to gather from
// the original context explicitly because the menu context becomes
// "current" once the popup opens.
func (h *ChordMenuHelper) gatherBindings() []*types.Binding {
	currentContext := h.chordStartContext()
	currentBindings := currentContext.GetKeybindings(h.c.KeybindingsOpts())
	globalBindings := h.c.Contexts().Global.GetKeybindings(h.c.KeybindingsOpts())
	return dedupBindingsByChord(currentBindings, globalBindings)
}

// dedupBindingsByChord returns currentBindings with non-duplicate
// globalBindings appended. Dedup keys on the canonical chord-label
// string rather than the gocui.Key value, because gocui.Key carries a
// *[]Key pointer for chord tails (see pkg/gocui/key.go) and two
// independently-parsed chord keys therefore have distinct identities
// even when they represent the same sequence.
//
// Internal duplicates within currentBindings are intentionally preserved:
// mode-toggling controllers (e.g. bisect) register two bindings for the
// same physical key, distinguished by AllowFurtherDispatching on the
// inactive one. BuildChordContinuations filters the inactive one out;
// dedup must not pre-empt that decision by collapsing them here.
func dedupBindingsByChord(currentBindings, globalBindings []*types.Binding) []*types.Binding {
	keyOf := func(b *types.Binding) string {
		return config.LabelForKeySequence(b.Key.Sequence())
	}
	seen := map[string]struct{}{}
	for _, b := range currentBindings {
		seen[keyOf(b)] = struct{}{}
	}
	for _, b := range globalBindings {
		if _, dup := seen[keyOf(b)]; dup {
			continue
		}
		currentBindings = append(currentBindings, b)
	}
	return currentBindings
}

// chordStartContext prefers originContext (which survives Pop+re-Push
// during group continuations), then gocui's pendingChordView, then the
// current context.
//
// If the current context is the Menu context we fall back to Global
// instead — the Menu's own bindings (j/k/Enter/Esc/etc.) are not
// useful as chord-continuation candidates, and pinning here would
// also poison originContext on the next openMenu pass.
func (h *ChordMenuHelper) chordStartContext() types.Context {
	h.mu.Lock()
	originContext := h.originContext
	h.mu.Unlock()
	if originContext != nil {
		return originContext
	}
	viewName := h.c.GocuiGui().PendingChordView()
	if viewName != "" {
		if ctx, ok := h.viewHelper.ContextForView(viewName); ok {
			return ctx
		}
	}
	current := h.c.Context().Current()
	if current.GetKey() == h.c.Contexts().Menu.GetKey() {
		h.c.Log.Warn("ChordMenuHelper.chordStartContext: current is Menu; falling back to Global")
		return h.c.Contexts().Global
	}
	return current
}

func (h *ChordMenuHelper) contextNameForChord() string {
	return string(h.chordStartContext().GetKey())
}
