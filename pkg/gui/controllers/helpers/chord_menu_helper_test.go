package helpers

import (
	"io"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// captureHook is a minimal logrus.Hook that records every entry. We use
// this rather than logrus/hooks/test because that subpackage isn't
// vendored.
type captureHook struct{ entries []*logrus.Entry }

func (h *captureHook) Levels() []logrus.Level { return logrus.AllLevels }
func (h *captureHook) Fire(e *logrus.Entry) error {
	h.entries = append(h.entries, e)
	return nil
}

// newBuildMenuItemsHelper returns a helper wired up enough for the
// buildMenuItems tests: just a logger entry. The returned hook captures
// any log writes the call emits.
func newBuildMenuItemsHelper(t *testing.T) (*ChordMenuHelper, *captureHook) {
	t.Helper()
	logger := logrus.New()
	logger.Out = io.Discard
	hook := &captureHook{}
	logger.AddHook(hook)
	c := &common.Common{Log: logrus.NewEntry(logger)}
	return &ChordMenuHelper{c: &HelperCommon{Common: c}}, hook
}

func TestTitleForPrefix(t *testing.T) {
	groups := map[string]map[string]config.KeybindingGroupConfig{
		"files":    {"i": {Name: "Ignore options"}},
		"branches": {"i": {Name: "Git flow options"}},
		"global":   {"m": {Name: "Rebase options"}},
	}
	tests := []struct {
		name     string
		prefix   []gocui.Key
		ctxName  string
		expected string
	}{
		{"context-specific wins", []gocui.Key{gocui.NewKeyRune('i')}, "files", "Ignore options"},
		{"different context different name", []gocui.Key{gocui.NewKeyRune('i')}, "branches", "Git flow options"},
		{"global fallback when ctx-specific absent", []gocui.Key{gocui.NewKeyRune('m')}, "files", "Rebase options"},
		{"generic fallback when neither matches", []gocui.Key{gocui.NewKeyRune('z')}, "files", "Chord: z …"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &ChordMenuHelper{}
			assert.Equal(t, tc.expected, h.titleForPrefix(tc.prefix, tc.ctxName, groups))
		})
	}
}

func TestBuildMenuItems_LeafRow(t *testing.T) {
	called := false
	binding := &types.Binding{
		Tooltip:         "Mark this commit as bad in the bisect",
		ChordPopupExtra: "git bisect bad <hash>",
		Handler:         func() error { called = true; return nil },
	}
	prefix := []gocui.Key{gocui.NewKeyRune('b')}

	infos := []bindingInfo{{
		key:         "b",
		description: "Mark as bad",
		tooltip:     "git bisect bad <hash>",
		binding:     binding,
	}}
	h, _ := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, prefix, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1)
	assert.Equal(t, []string{"Mark as bad", "git bisect bad <hash>"}, items[0].LabelColumns)
	assert.Equal(t, gocui.NewKeyRune('b'), items[0].Key)
	assert.Equal(t, "Mark this commit as bad in the bisect", items[0].Tooltip)
	assert.NoError(t, items[0].OnPress())
	assert.True(t, called, "binding handler should have run")
}

func TestBuildMenuItems_GroupRow(t *testing.T) {
	prefix := []gocui.Key{}
	infos := []bindingInfo{{
		key:         "M",
		description: "Merge",
		isGroup:     true,
	}}

	var extendedTo []gocui.Key
	h, _ := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, prefix, func(p []gocui.Key) error {
		extendedTo = p
		return nil
	})

	assert.Len(t, items, 1)
	assert.Equal(t, gocui.NewKeyRune('M'), items[0].Key)
	assert.NoError(t, items[0].OnPress())
	assert.Equal(t, []gocui.Key{gocui.NewKeyRune('M')}, extendedTo)
}

func TestBuildMenuItems_LeafRowMirrorsDisabledReason(t *testing.T) {
	reason := &types.DisabledReason{Text: "no commit selected"}
	binding := &types.Binding{
		Handler:           func() error { return nil },
		GetDisabledReason: func() *types.DisabledReason { return reason },
	}
	infos := []bindingInfo{{key: "b", description: "Mark as bad", binding: binding}}

	h, _ := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1)
	assert.Same(t, reason, items[0].DisabledReason)
}

func TestBuildMenuItems_LeafRowEnabledHasNilReason(t *testing.T) {
	binding := &types.Binding{Handler: func() error { return nil }}
	infos := []bindingInfo{{key: "b", description: "Mark as bad", binding: binding}}

	h, _ := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1)
	assert.Nil(t, items[0].DisabledReason)
}

func TestBuildMenuItems_AlignsExtraColumn(t *testing.T) {
	infos := []bindingInfo{
		{key: "a", description: "Alpha", tooltip: "", binding: &types.Binding{Handler: func() error { return nil }}},
		{key: "b", description: "Beta", tooltip: "git beta", binding: &types.Binding{Handler: func() error { return nil }}},
	}
	h, _ := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items[0].LabelColumns, 2)
	assert.Equal(t, "", items[0].LabelColumns[1])
	assert.Equal(t, "git beta", items[1].LabelColumns[1])
}

func TestBuildMenuItems_LogsAndSkipsUnparseableKey(t *testing.T) {
	// Sanity-check the chosen label actually fails KeyFromLabel; if it
	// ever starts parsing the test needs a different unparseable label.
	if _, ok := config.KeyFromLabel("<<<"); ok {
		t.Fatalf(`expected KeyFromLabel("<<<") to fail`)
	}

	infos := []bindingInfo{
		{key: "<<<", description: "broken", binding: &types.Binding{Handler: func() error { return nil }}},
		{key: "b", description: "Mark as bad", binding: &types.Binding{Handler: func() error { return nil }}},
	}

	h, hook := newBuildMenuItemsHelper(t)
	items := h.buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1, "malformed row must be skipped")
	assert.Equal(t, gocui.NewKeyRune('b'), items[0].Key)

	assert.Len(t, hook.entries, 1, "exactly one error log expected")
	assert.Equal(t, logrus.ErrorLevel, hook.entries[0].Level)
	assert.Contains(t, hook.entries[0].Message, "<<<")
}

func TestBuildChordContinuations_ContextSpecificGroupName(t *testing.T) {
	// "i" is a group prefix in both files and branches; the helper
	// must resolve based on ctxName.
	groups := map[string]map[string]config.KeybindingGroupConfig{
		"files":    {"i": {Name: "Ignore options"}},
		"branches": {"i": {Name: "Git flow options"}},
	}
	keyII, _ := config.KeyFromLabel("ii")
	keyIE, _ := config.KeyFromLabel("ie")
	bindings := []*types.Binding{
		{Key: keyII, Description: "Add to gitignore"},
		{Key: keyIE, Description: "Add to git/info/exclude"},
	}
	rows := BuildChordContinuations(bindings, nil, groups, "files")
	assert.Len(t, rows, 1)
	assert.Equal(t, "i", rows[0].key)
	assert.Equal(t, "Ignore options", rows[0].description)
	assert.True(t, rows[0].isGroup)
}

type fakeJob struct {
	fn        func()
	cancelled bool
}

func (j *fakeJob) Run() {
	if !j.cancelled {
		j.fn()
	}
}

// ForceRun fires the closure regardless of cancellation. Models the
// time.AfterFunc race where the callback was already queued on the UI
// thread when Stop() was called — ChordMenuHelper must defend against
// this with the timerGen counter.
func (j *fakeJob) ForceRun() { j.fn() }

// Tests advance time by calling f.scheduled[i].Run() explicitly.
type fakeScheduler struct {
	scheduled []*fakeJob
}

func (f *fakeScheduler) Schedule(_ time.Duration, fn func()) func() {
	job := &fakeJob{fn: fn}
	f.scheduled = append(f.scheduled, job)
	return func() { job.cancelled = true }
}

// UserConfig() is provided by the embedded *common.Common (which holds
// the config in an atomic.Pointer), shadowing the IGuiCommon method on
// the embedded struct.
func newTestHelperCommon(cfg *config.UserConfig) *HelperCommon {
	c := &common.Common{}
	c.SetUserConfig(cfg)
	return &HelperCommon{Common: c}
}

func newTestHelper(t *testing.T, delayMs int) *ChordMenuHelper {
	t.Helper()
	cfg := &config.UserConfig{ChordPopupDelayMs: delayMs}
	return &ChordMenuHelper{c: newTestHelperCommon(cfg), scheduler: realScheduler{}}
}

func newTestHelperWithScheduler(t *testing.T, delayMs int, sched scheduler) *ChordMenuHelper {
	t.Helper()
	cfg := &config.UserConfig{ChordPopupDelayMs: delayMs}
	return &ChordMenuHelper{c: newTestHelperCommon(cfg), scheduler: sched}
}

func TestChordMenuHelper_NoOpenWhenDelayDisabled(t *testing.T) {
	h := newTestHelper(t, -1)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Equal(t, 0, opens)
}

func TestChordMenuHelper_InstantOpen(t *testing.T) {
	h := newTestHelper(t, 0)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Equal(t, 1, opens)
	assert.True(t, h.menuOpen)
}

func TestChordMenuHelper_DelayedOpen(t *testing.T) {
	sched := &fakeScheduler{}
	h := newTestHelperWithScheduler(t, 200, sched)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Equal(t, 0, opens)
	assert.Len(t, sched.scheduled, 1)
	sched.scheduled[0].Run()
	assert.Equal(t, 1, opens)
}

func TestChordMenuHelper_EmptyPrefixCancelsTimer(t *testing.T) {
	sched := &fakeScheduler{}
	h := newTestHelperWithScheduler(t, 200, sched)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	h.OnChordStateChange(nil)
	assert.Len(t, sched.scheduled, 1)
	sched.scheduled[0].Run()
	assert.Equal(t, 0, opens, "open should not happen after cancel")
}

func TestChordMenuHelper_IgnoresReentrantEmptyDuringOpen(t *testing.T) {
	h := newTestHelper(t, 0)
	opens := 0
	closes := 0
	h.openHookForTest = func(_ []gocui.Key) {
		opens++
		// Simulate the re-entrant empty callback that Push triggers.
		h.OnChordStateChange(nil)
	}
	h.closeHookForTest = func() { closes++ }
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Equal(t, 1, opens)
	assert.Equal(t, 0, closes, "re-entrant empty callback during open must be ignored")
}

func TestChordMenuHelper_NotifyMenuClosedResetsState(t *testing.T) {
	h := newTestHelper(t, 0)
	h.openHookForTest = func(_ []gocui.Key) {}
	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.True(t, h.menuOpen)
	h.NotifyMenuClosed()
	assert.False(t, h.menuOpen)
}

func TestChordMenuHelper_TitleFunc_TakesPrecedenceOverGroupName(t *testing.T) {
	h := &ChordMenuHelper{}
	h.RegisterTitleFunc("localBranches", "M", func() string { return "Merge into 'main'" })

	groups := map[string]map[string]config.KeybindingGroupConfig{
		"localBranches": {"M": {Name: "Merge"}},
	}
	keyM, _ := config.KeyFromLabel("M")

	title := h.titleForPrefix(keyM.Sequence(), "localBranches", groups)
	assert.Equal(t, "Merge into 'main'", title)
}

func TestChordMenuHelper_TitleFunc_FallsBackToGroupName(t *testing.T) {
	h := &ChordMenuHelper{}
	groups := map[string]map[string]config.KeybindingGroupConfig{
		"localBranches": {"M": {Name: "Merge"}},
	}
	keyM, _ := config.KeyFromLabel("M")

	title := h.titleForPrefix(keyM.Sequence(), "localBranches", groups)
	assert.Equal(t, "Merge", title)
}

func TestChordMenuHelper_TitleFunc_GlobalFallback(t *testing.T) {
	h := &ChordMenuHelper{}
	h.RegisterTitleFunc("global", "X", func() string { return "Global X" })

	groups := map[string]map[string]config.KeybindingGroupConfig{}
	keyX, _ := config.KeyFromLabel("X")

	title := h.titleForPrefix(keyX.Sequence(), "files", groups)
	assert.Equal(t, "Global X", title)
}

func TestGatherBindings_DedupesChordSharedAcrossContexts(t *testing.T) {
	// Parse the same chord label twice. Each call returns a fresh
	// gocui.Key whose `rest` field is a distinct *[]Key pointer, so the
	// two values are NOT equal as map keys even though they represent
	// the same chord sequence.
	keyA, okA := config.KeyFromLabel("<b><p>")
	keyB, okB := config.KeyFromLabel("<b><p>")
	assert.True(t, okA)
	assert.True(t, okB)

	currentBindings := []*types.Binding{{Key: keyA, Description: "current"}}
	globalBindings := []*types.Binding{{Key: keyB, Description: "global"}}

	deduped := dedupBindingsByChord(currentBindings, globalBindings)

	assert.Len(t, deduped, 1, "chord bindings parsed independently must dedupe by canonical label")
	assert.Equal(t, "current", deduped[0].Description, "current-context binding must take precedence")
}

func TestGatherBindings_KeepsDistinctChords(t *testing.T) {
	keyBP, _ := config.KeyFromLabel("<b><p>")
	keyBN, _ := config.KeyFromLabel("<b><n>")

	currentBindings := []*types.Binding{{Key: keyBP, Description: "current bp"}}
	globalBindings := []*types.Binding{{Key: keyBN, Description: "global bn"}}

	deduped := dedupBindingsByChord(currentBindings, globalBindings)

	assert.Len(t, deduped, 2, "distinct chord sequences must not be deduped")
}

// Models the documented race: the timer was Stop()'d but its closure
// was already queued for execution. The fake job's ForceRun bypasses
// the cancelled flag to emulate the in-flight callback. The helper's
// timerGen counter must detect that the chord state has moved on and
// suppress the open.
func TestChordMenuHelper_TimerFiresAfterCancel(t *testing.T) {
	sched := &fakeScheduler{}
	h := newTestHelperWithScheduler(t, 200, sched)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }

	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Len(t, sched.scheduled, 1)

	h.OnChordStateChange(nil)

	sched.scheduled[0].ForceRun()
	assert.Equal(t, 0, opens, "in-flight timer fired after cancel must not open")
}

// Same race but the chord was *replaced* rather than cancelled — a
// new prefix bumped the generation. The original timer's queued
// callback must still bail.
func TestChordMenuHelper_TimerFiresAfterReplacement(t *testing.T) {
	sched := &fakeScheduler{}
	h := newTestHelperWithScheduler(t, 200, sched)
	opens := 0
	h.openHookForTest = func(_ []gocui.Key) { opens++ }

	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})
	assert.Len(t, sched.scheduled, 1)

	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('m')})
	assert.Len(t, sched.scheduled, 2)

	sched.scheduled[0].ForceRun()
	assert.Equal(t, 0, opens, "first timer must bail after generation bump")

	sched.scheduled[1].Run()
	assert.Equal(t, 1, opens, "second (current-generation) timer must open")
}

// The opening guard must drop only the empty re-entrant Push
// callback. Non-empty changes arriving mid-open must be queued for
// replay after openMenu completes — running them inline would race
// with the in-flight c.Menu call.
func TestChordMenuHelper_NonEmptyCallbackDuringOpenIsDeferred(t *testing.T) {
	h := newTestHelper(t, 0)
	deferredObservedDuringHook := false
	h.openHookForTest = func(_ []gocui.Key) {
		h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b'), gocui.NewKeyRune('p')})

		h.mu.Lock()
		deferredObservedDuringHook = h.deferredPrefix != nil
		h.mu.Unlock()
	}

	var replayedTo [][]gocui.Key
	h.refreshHookForTest = func(p []gocui.Key) {
		replayedTo = append(replayedTo, append([]gocui.Key(nil), p...))
	}

	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})

	assert.True(t, deferredObservedDuringHook,
		"non-empty change during opening must be captured as deferred, not processed inline")
	assert.Len(t, replayedTo, 1, "deferred prefix must be replayed after open completes")
	assert.Equal(t, []gocui.Key{gocui.NewKeyRune('b'), gocui.NewKeyRune('p')}, replayedTo[0])
}

// Confirms the documented re-entrant empty callback is dropped, not
// deferred — deferring would replay an empty prefix and incorrectly
// close the menu we just pushed.
func TestChordMenuHelper_EmptyCallbackDuringOpenIsDroppedNotDeferred(t *testing.T) {
	h := newTestHelper(t, 0)
	h.openHookForTest = func(_ []gocui.Key) {
		h.OnChordStateChange(nil)
	}

	h.OnChordStateChange([]gocui.Key{gocui.NewKeyRune('b')})

	h.mu.Lock()
	hasDeferred := h.deferredPrefix != nil
	h.mu.Unlock()

	assert.False(t, hasDeferred,
		"empty re-entrant Push callback must be dropped, not deferred")
}
