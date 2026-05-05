package config

import "testing"

func TestChordPrefixLookup_Defaults(t *testing.T) {
	cfg := GetDefaultConfig()
	cp := cfg.Keybinding.ChordPrefix

	cases := []struct {
		name    string
		context string
		id      string
		want    string
	}{
		{"global rebase options", "global", ChordIDRebaseOptions, "m"},
		{"files ignore options", "files", ChordIDIgnoreOptions, "i"},
		{"files stash options", "files", ChordIDStashOptions, "S"},
		{"files discard/reset options", "files", ChordIDDiscardAndResetOptions, "D"},
		{"files filter", "files", ChordIDFilterFiles, "<ctrl+b>"},
		{"files copy to clipboard", "files", ChordIDCopyToClipboard, "y"},
		{"files reset to upstream", "files", ChordIDResetToUpstream, "g"},
		{"files discard changes", "files", ChordIDDiscardChanges, "d"},
		{"localBranches merge", "localBranches", ChordIDMerge, "M"},
		{"localBranches delete branch", "localBranches", ChordIDDeleteBranch, "d"},
		{"localBranches git flow", "localBranches", ChordIDGitFlowOptions, "i"},
		{"localBranches rebase options", "localBranches", ChordIDRebaseOptions, "r"},
		{"localBranches branch upstream", "localBranches", ChordIDBranchUpstreamOptions, "u"},
		{"localBranches reset to ref", "localBranches", ChordIDResetToRef, "g"},
		{"remoteBranches delete remote", "remoteBranches", ChordIDDeleteRemoteBranch, "d"},
		{"remoteBranches merge", "remoteBranches", ChordIDMerge, "M"},
		{"tags delete tag", "tags", ChordIDDeleteTag, "d"},
		{"commits bisect", "commits", ChordIDBisectOptions, "b"},
		{"commits fixup", "commits", ChordIDFixupCommitOptions, "f"},
		{"commits reset to ref", "commits", ChordIDResetToRef, "g"},
		{"commitFiles copy to clipboard", "commitFiles", ChordIDCopyToClipboard, "y"},
		{"submodules bulk options", "submodules", ChordIDBulkOptions, "b"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := cp.Get(c.context, c.id); got != c.want {
				t.Errorf("Get(%q, %q) = %q, want %q", c.context, c.id, got, c.want)
			}
		})
	}
}

func TestChordPrefixLookup_FallsBackToGlobal(t *testing.T) {
	cfg := GetDefaultConfig()
	cp := cfg.Keybinding.ChordPrefix

	// "global" is the only context that registers ChordIDRebaseOptions
	// alongside its localBranches counterpart. A context that doesn't
	// register the ID at all should fall back to global.
	if got := cp.Get("worktrees", ChordIDRebaseOptions); got != "m" {
		t.Errorf("fallback to global: got %q, want %q", got, "m")
	}
}

func TestChordPrefixLookup_PrefersContextOverGlobal(t *testing.T) {
	cfg := GetDefaultConfig()
	cp := cfg.Keybinding.ChordPrefix

	// localBranches has its own RebaseOptions ("r"); global has "m".
	// The contextual one wins.
	if got := cp.Get("localBranches", ChordIDRebaseOptions); got != "r" {
		t.Errorf("contextual preference: got %q, want %q", got, "r")
	}
}

func TestChordPrefixLookup_MissingReturnsEmpty(t *testing.T) {
	cfg := GetDefaultConfig()
	cp := cfg.Keybinding.ChordPrefix

	if got := cp.Get("commits", "noSuchID"); got != "" {
		t.Errorf("unknown id: got %q, want empty", got)
	}
	if got := cp.Get("commits", ""); got != "" {
		t.Errorf("empty id: got %q, want empty", got)
	}
}

func TestChordPrefixLookup_RespondsToCustomization(t *testing.T) {
	cfg := GetDefaultConfig()
	// Move the global rebase chord from "m" to "x" by re-keying.
	cfg.KeybindingGroups["global"] = map[string]KeybindingGroupConfig{
		"x": {ID: ChordIDRebaseOptions, Name: "Rebase options"},
	}
	cfg.ResolveChordPrefixes()

	if got := cfg.Keybinding.ChordPrefix.Get("global", ChordIDRebaseOptions); got != "x" {
		t.Errorf("after remap: got %q, want %q", got, "x")
	}
}
