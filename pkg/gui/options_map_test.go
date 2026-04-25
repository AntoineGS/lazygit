package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func mustKey(t *testing.T, label string) gocui.Key {
	t.Helper()
	k, ok := config.KeyFromLabel(label)
	if !ok {
		t.Fatalf("KeyFromLabel(%q) failed", label)
	}
	return k
}

func TestBuildChordContinuations_LeavesAndGroupsCollapse(t *testing.T) {
	bindings := []*types.Binding{
		{Key: mustKey(t, "<b><p>"), Description: "Pull"},
		{Key: mustKey(t, "<b><P>"), Description: "Push"},
		{Key: mustKey(t, "<b><t><o>"), Description: "Open PR"},
		{Key: mustKey(t, "<b><t><l>"), Description: "List PRs"},
	}
	prefix := mustKey(t, "<b>").Sequence()
	groups := map[string]config.KeybindingGroupConfig{
		"<b><t>": {Name: "Pull Request"},
	}

	got := buildChordContinuations(bindings, prefix, groups)
	// Expect: p:Pull, P:Push, t:Pull Request, <esc>:cancel — 4 rows.
	if len(got) != 4 {
		t.Fatalf("expected 4 rows, got %d: %+v", len(got), got)
	}
	descByKey := map[string]string{}
	for _, r := range got {
		descByKey[r.key] = r.description
	}
	if descByKey["t"] != "Pull Request" {
		t.Errorf("expected t row to use group name 'Pull Request', got %q", descByKey["t"])
	}
	if descByKey["p"] != "Pull" {
		t.Errorf("expected p row to use binding description, got %q", descByKey["p"])
	}
	if _, ok := descByKey["<esc>"]; !ok {
		t.Errorf("missing <esc> cancel row")
	}
}
