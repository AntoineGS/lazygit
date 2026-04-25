package helpers

import (
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, tc.expected, titleForPrefix(tc.prefix, tc.ctxName, groups))
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
	items := buildMenuItems(infos, prefix, func(_ []gocui.Key) error { return nil })

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
	items := buildMenuItems(infos, prefix, func(p []gocui.Key) error {
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

	items := buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1)
	assert.Same(t, reason, items[0].DisabledReason)
}

func TestBuildMenuItems_LeafRowEnabledHasNilReason(t *testing.T) {
	binding := &types.Binding{Handler: func() error { return nil }}
	infos := []bindingInfo{{key: "b", description: "Mark as bad", binding: binding}}

	items := buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items, 1)
	assert.Nil(t, items[0].DisabledReason)
}

func TestBuildMenuItems_AlignsExtraColumn(t *testing.T) {
	infos := []bindingInfo{
		{key: "a", description: "Alpha", tooltip: "", binding: &types.Binding{Handler: func() error { return nil }}},
		{key: "b", description: "Beta", tooltip: "git beta", binding: &types.Binding{Handler: func() error { return nil }}},
	}
	items := buildMenuItems(infos, nil, func(_ []gocui.Key) error { return nil })

	assert.Len(t, items[0].LabelColumns, 2)
	assert.Equal(t, "", items[0].LabelColumns[1])
	assert.Equal(t, "git beta", items[1].LabelColumns[1])
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
	assert.Nil(t, h.openPrefix)
}
