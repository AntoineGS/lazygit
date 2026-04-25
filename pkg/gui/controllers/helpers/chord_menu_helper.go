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
	openPrefix  []gocui.Key
	menuOpen    bool
	opening     bool
	// originContext is preserved across the menu Push/Pop so continuations
	// remain scoped to the originating view even though gocui's chord
	// state is cleared by Push.
	originContext types.Context

	openHookForTest  func([]gocui.Key)
	closeHookForTest func()
}

func NewChordMenuHelper(c *HelperCommon, viewHelper *ViewHelper) *ChordMenuHelper {
	return &ChordMenuHelper{c: c, viewHelper: viewHelper, scheduler: realScheduler{}}
}

func titleForPrefix(
	prefix []gocui.Key,
	ctxName string,
	groups map[string]map[string]config.KeybindingGroupConfig,
) string {
	label := config.LabelForKeySequence(prefix)
	if g, ok := groups[ctxName][label]; ok && g.Name != "" {
		return g.Name
	}
	if g, ok := groups["global"][label]; ok && g.Name != "" {
		return g.Name
	}
	return fmt.Sprintf("Chord: %s …", label)
}

// buildMenuItems maps chord-continuation rows into MenuItems.
//
// If any row has a non-empty inline tooltip, every row gets a 2nd
// LabelColumns entry — possibly empty — so column widths align.
func buildMenuItems(
	infos []bindingInfo,
	prefix []gocui.Key,
	extendFn func([]gocui.Key) error,
) []*types.MenuItem {
	anyExtra := lo.SomeBy(infos, func(i bindingInfo) bool { return i.tooltip != "" })

	return lo.Map(infos, func(info bindingInfo, _ int) *types.MenuItem {
		cols := []string{info.description}
		if anyExtra {
			cols = append(cols, info.tooltip)
		}

		nextKey, _ := config.KeyFromLabel(info.key)
		nextSeq := nextKey.Sequence()
		var thisKey gocui.Key
		if len(nextSeq) > 0 {
			thisKey = nextSeq[0]
		}

		item := &types.MenuItem{
			LabelColumns: cols,
			Key:          thisKey,
		}

		if info.isGroup {
			extended := append([]gocui.Key{}, prefix...)
			extended = append(extended, thisKey)
			item.OnPress = func() error { return extendFn(extended) }
			return item
		}

		// Mirror GetDisabledReason onto the MenuItem so disabled rows
		// surface the standard toast on press, matching the regular
		// dispatch path. AllowFurtherDispatching rows are filtered out
		// upstream by BuildChordContinuations.
		if info.binding.GetDisabledReason != nil {
			item.DisabledReason = info.binding.GetDisabledReason()
		}
		item.OnPress = info.binding.Handler
		item.Tooltip = info.binding.Tooltip
		return item
	})
}

// OnChordStateChange is invoked by gocui when the chord state changes.
// A negative ChordPopupDelayMs disables the popup entirely.
//
// The opening flag guards against the re-entrant empty callback that
// c.Menu's Push triggers while we're opening our own menu.
func (h *ChordMenuHelper) OnChordStateChange(prefix []gocui.Key) {
	h.mu.Lock()
	if h.opening {
		h.mu.Unlock()
		return
	}
	delayMs := h.c.UserConfig().ChordPopupDelayMs
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
		h.openPrefix = append([]gocui.Key(nil), prefix...)
		h.menuOpen = true
		h.mu.Unlock()
		h.openHookForTest(prefix)
		h.mu.Lock()
		h.opening = false
		h.mu.Unlock()
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
	items := buildMenuItems(infos, prefix, h.extendPrefix)

	h.mu.Lock()
	h.opening = true
	h.openPrefix = append([]gocui.Key(nil), prefix...)
	h.mu.Unlock()

	err := h.c.Menu(types.CreateMenuOptions{
		Title: titleForPrefix(prefix, ctxName, groups),
		Items: items,
	})

	h.mu.Lock()
	h.opening = false
	if err == nil {
		h.menuOpen = true
	}
	h.mu.Unlock()
}

// refreshMenu swaps the open menu's items and title in place — no Push
// or Pop. We re-register keybindings so the new continuation keys
// dispatch via the menu's per-item bindings.
func (h *ChordMenuHelper) refreshMenu(newPrefix []gocui.Key) {
	groups := h.c.UserConfig().KeybindingGroups
	ctxName := h.contextNameForChord()
	infos := BuildChordContinuations(h.gatherBindings(), newPrefix, groups, ctxName)
	items := buildMenuItems(infos, newPrefix, h.extendPrefix)

	h.c.Contexts().Menu.SetMenuItems(items, nil)
	h.c.Views().Menu.Title = titleForPrefix(newPrefix, ctxName, groups)
	if err := h.c.ResetKeybindings(); err != nil {
		h.c.Log.Errorf("ChordMenu refreshMenu ResetKeybindings: %v", err)
	}
	h.c.PostRefreshUpdate(h.c.Contexts().Menu)

	h.mu.Lock()
	h.openPrefix = append([]gocui.Key(nil), newPrefix...)
	h.mu.Unlock()
}

func (h *ChordMenuHelper) closeMenu() {
	h.mu.Lock()
	if !h.menuOpen {
		h.mu.Unlock()
		return
	}
	h.menuOpen = false
	h.openPrefix = nil
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
	h.openPrefix = nil
	h.mu.Unlock()
	h.openMenu(newPrefix)
	return nil
}

func (h *ChordMenuHelper) NotifyMenuClosed() {
	h.mu.Lock()
	h.menuOpen = false
	h.openPrefix = nil
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
	seen := map[gocui.Key]struct{}{}
	for _, b := range currentBindings {
		seen[b.Key] = struct{}{}
	}
	for _, b := range globalBindings {
		if _, dup := seen[b.Key]; dup {
			continue
		}
		currentBindings = append(currentBindings, b)
	}
	return currentBindings
}

// chordStartContext prefers originContext (which survives Pop+re-Push
// during group continuations), then gocui's pendingChordView, then the
// current context.
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
	return h.c.Context().Current()
}

func (h *ChordMenuHelper) contextNameForChord() string {
	return string(h.chordStartContext().GetKey())
}
