package config

import "testing"

func TestResolveChordPrefixes_PopulatesFromDefaultGroups(t *testing.T) {
	cfg := GetDefaultConfig()
	cp := cfg.Keybinding.ChordPrefix

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"Global.RebaseOptions", cp.Global.RebaseOptions, "m"},
		{"Files.IgnoreOptions", cp.Files.IgnoreOptions, "i"},
		{"Files.StashOptions", cp.Files.StashOptions, "S"},
		{"Files.DiscardAndResetOptions", cp.Files.DiscardAndResetOptions, "D"},
		{"Files.FilterFiles", cp.Files.FilterFiles, "<ctrl+b>"},
		{"Files.CopyToClipboard", cp.Files.CopyToClipboard, "y"},
		{"LocalBranches.Merge", cp.LocalBranches.Merge, "M"},
		{"LocalBranches.DeleteBranch", cp.LocalBranches.DeleteBranch, "d"},
		{"LocalBranches.GitFlowOptions", cp.LocalBranches.GitFlowOptions, "i"},
		{"LocalBranches.RebaseOptions", cp.LocalBranches.RebaseOptions, "r"},
		{"LocalBranches.BranchUpstreamOptions", cp.LocalBranches.BranchUpstreamOptions, "u"},
		{"RemoteBranches.DeleteRemoteBranch", cp.RemoteBranches.DeleteRemoteBranch, "d"},
		{"Tags.DeleteTag", cp.Tags.DeleteTag, "d"},
		{"Commits.BisectOptions", cp.Commits.BisectOptions, "b"},
		{"Commits.FixupCommitOptions", cp.Commits.FixupCommitOptions, "f"},
		{"Commits.ResetToRef", cp.Commits.ResetToRef, "g"},
		{"Submodules.BulkOptions", cp.Submodules.BulkOptions, "b"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestResolveChordPrefixes_LeavesUnknownGroupsEmpty(t *testing.T) {
	cfg := &UserConfig{
		KeybindingGroups: map[string]map[string]KeybindingGroupConfig{
			"global": {
				"X": {Name: "Some custom group"},
			},
		},
	}
	cfg.ResolveChordPrefixes()

	if cfg.Keybinding.ChordPrefix.Global.RebaseOptions != "" {
		t.Errorf("expected empty when no matching group, got %q",
			cfg.Keybinding.ChordPrefix.Global.RebaseOptions)
	}
}

func TestResolveChordPrefixes_RespondsToCustomization(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.KeybindingGroups["global"] = map[string]KeybindingGroupConfig{
		"x": {Name: "Rebase options"},
	}
	cfg.ResolveChordPrefixes()

	if got := cfg.Keybinding.ChordPrefix.Global.RebaseOptions; got != "x" {
		t.Errorf("after remap, Global.RebaseOptions = %q, want %q", got, "x")
	}
}
