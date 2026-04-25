package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserConfigValidate_enums(t *testing.T) {
	type testCase struct {
		value string
		valid bool
	}

	scenarios := []struct {
		name      string
		setup     func(config *UserConfig, value string)
		testCases []testCase
	}{
		{
			name: "Gui.StatusPanelView",
			setup: func(config *UserConfig, value string) {
				config.Gui.StatusPanelView = value
			},
			testCases: []testCase{
				{value: "dashboard", valid: true},
				{value: "allBranchesLog", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Gui.ShowDivergenceFromBaseBranch",
			setup: func(config *UserConfig, value string) {
				config.Gui.ShowDivergenceFromBaseBranch = value
			},
			testCases: []testCase{
				{value: "none", valid: true},
				{value: "onlyArrow", valid: true},
				{value: "arrowAndNumber", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.AutoForwardBranches",
			setup: func(config *UserConfig, value string) {
				config.Git.AutoForwardBranches = value
			},
			testCases: []testCase{
				{value: "none", valid: true},
				{value: "onlyMainBranches", valid: true},
				{value: "allBranches", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.LocalBranchSortOrder",
			setup: func(config *UserConfig, value string) {
				config.Git.LocalBranchSortOrder = value
			},
			testCases: []testCase{
				{value: "date", valid: true},
				{value: "recency", valid: true},
				{value: "alphabetical", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.RemoteBranchSortOrder",
			setup: func(config *UserConfig, value string) {
				config.Git.RemoteBranchSortOrder = value
			},
			testCases: []testCase{
				{value: "date", valid: true},
				{value: "recency", valid: false},
				{value: "alphabetical", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.Log.Order",
			setup: func(config *UserConfig, value string) {
				config.Git.Log.Order = value
			},
			testCases: []testCase{
				{value: "date-order", valid: true},
				{value: "author-date-order", valid: true},
				{value: "topo-order", valid: true},
				{value: "default", valid: true},

				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Git.Log.ShowGraph",
			setup: func(config *UserConfig, value string) {
				config.Git.Log.ShowGraph = value
			},
			testCases: []testCase{
				{value: "always", valid: true},
				{value: "never", valid: true},
				{value: "when-maximised", valid: true},

				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Keybindings",
			setup: func(config *UserConfig, value string) {
				config.Keybinding.Universal.Quit = value
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				// Multi-rune strings are valid as chord sequences;
				// unterminated bracket asserts a rejected key.
				{value: "<bogus", valid: false},
			},
		},
		{
			name: "JumpToBlock keybinding",
			setup: func(config *UserConfig, value string) {
				config.Keybinding.Universal.JumpToBlock = strings.Split(value, ",")
			},
			testCases: []testCase{
				{value: "", valid: false},
				{value: "1,2,3", valid: false},
				{value: "1,2,3,4,5", valid: true},
				{value: "1,2,3,4,<bogus", valid: false},
				{value: "1,2,3,4,5,6", valid: false},
			},
		},
		{
			name: "Custom command keybinding",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     value,
						Command: "echo 'hello'",
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "<bogus", valid: false},
			},
		},
		{
			name: "Custom command keybinding in sub menu",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         "X",
						Description: "My Custom Commands",
						CommandMenu: []CustomCommand{
							{Key: value, Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "<bogus", valid: false},
			},
		},
		{
			name: "Custom command keybinding in prompt menu",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         "X",
						Description: "My Custom Commands",
						Prompts: []CustomCommandPrompt{
							{
								Options: []CustomCommandMenuOption{
									{Key: value},
								},
							},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "<bogus", valid: false},
			},
		},
		{
			name: "Custom command output",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Output: value,
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "none", valid: true},
				{value: "terminal", valid: true},
				{value: "log", valid: true},
				{value: "logWithPty", valid: true},
				{value: "popup", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         "X",
						Description: "My Custom Commands",
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     "X",
						Context: "global", // context is not allowed for submenus
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         "X",
						LoadingText: "loading", // other properties are not allowed for submenus (using loadingText as an example)
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			for _, testCase := range s.testCases {
				config := GetDefaultConfig()
				s.setup(config, testCase.value)
				err := config.Validate()

				if testCase.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestKeybindingGroup_PrefixMustParse(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.KeybindingGroups = map[string]map[string]KeybindingGroupConfig{
		"global": {
			"<bogus": {Name: "Bad"},
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for unparsable prefix")
	}
	if !strings.Contains(err.Error(), "<bogus") {
		t.Fatalf("error should cite the offending prefix, got: %v", err)
	}
}

func TestKeybindingGroup_NameMustNotBeEmpty(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.Pull = "<X><p>"
	cfg.KeybindingGroups = map[string]map[string]KeybindingGroupConfig{
		"global": {
			"<X>": {Name: ""},
		},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "non-empty name") {
		t.Fatalf("expected validation error for empty name, got: %v", err)
	}
}

func TestKeybindingGroup_LeafCollisionRejected(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.Pull = "<b>"
	cfg.Keybinding.Universal.Push = "<b><p>"
	cfg.KeybindingGroups = map[string]map[string]KeybindingGroupConfig{
		"global": {
			"<b>": {Name: "Branch"},
		},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "<b>") {
		t.Fatalf("expected leaf/group collision error for <b>, got: %v", err)
	}
}

func TestKeybindingGroup_MustHaveAtLeastOneBinding(t *testing.T) {
	cfg := GetDefaultConfig()
	// No chord binding starts with <z>.
	cfg.KeybindingGroups = map[string]map[string]KeybindingGroupConfig{
		"global": {
			"<z>": {Name: "Empty"},
		},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "<z>") {
		t.Fatalf("expected error citing empty group <z>, got: %v", err)
	}
}

func TestChordPopupDelayMsDefault(t *testing.T) {
	cfg := GetDefaultConfig()
	if cfg.ChordPopupDelayMs != 0 {
		t.Fatalf("expected default ChordPopupDelayMs=0, got %d", cfg.ChordPopupDelayMs)
	}
}

func TestStashChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"StashAllChangesKeepIndex":     "Si",
		"StashIncludeUntrackedChanges": "SU",
		"StashStagedChanges":           "Ss",
		"StashUnstagedChanges":         "Su",
	}
	actual := map[string]string{
		"StashAllChangesKeepIndex":     cfg.Keybinding.Files.StashAllChangesKeepIndex,
		"StashIncludeUntrackedChanges": cfg.Keybinding.Files.StashIncludeUntrackedChanges,
		"StashStagedChanges":           cfg.Keybinding.Files.StashStagedChanges,
		"StashUnstagedChanges":         cfg.Keybinding.Files.StashUnstagedChanges,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestSubmoduleChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"BulkInit":            "bi",
		"BulkUpdate":          "bu",
		"BulkUpdateRecursive": "br",
		"BulkDeinit":          "bd",
	}
	actual := map[string]string{
		"BulkInit":            cfg.Keybinding.Submodules.BulkInit,
		"BulkUpdate":          cfg.Keybinding.Submodules.BulkUpdate,
		"BulkUpdateRecursive": cfg.Keybinding.Submodules.BulkUpdateRecursive,
		"BulkDeinit":          cfg.Keybinding.Submodules.BulkDeinit,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestFilesCopyChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"CopyFileName":         "yn",
		"CopyRelativeFilePath": "yp",
		"CopyAbsoluteFilePath": "yP",
		"CopyFileDiff":         "ys",
		"CopyAllFilesDiff":     "ya",
	}
	actual := map[string]string{
		"CopyFileName":         cfg.Keybinding.Files.CopyFileName,
		"CopyRelativeFilePath": cfg.Keybinding.Files.CopyRelativeFilePath,
		"CopyAbsoluteFilePath": cfg.Keybinding.Files.CopyAbsoluteFilePath,
		"CopyFileDiff":         cfg.Keybinding.Files.CopyFileDiff,
		"CopyAllFilesDiff":     cfg.Keybinding.Files.CopyAllFilesDiff,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestCommitFilesCopyContentChordDefault(t *testing.T) {
	cfg := GetDefaultConfig()
	if got := cfg.Keybinding.CommitFiles.CopyFileContent; got != "yc" {
		t.Errorf("default CopyFileContent = %q, want %q", got, "yc")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestFilesIgnoreExcludeChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	if got := cfg.Keybinding.Files.Ignore; got != "ii" {
		t.Errorf("default Ignore = %q, want %q", got, "ii")
	}
	if got := cfg.Keybinding.Files.Exclude; got != "ie" {
		t.Errorf("default Exclude = %q, want %q", got, "ie")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestResetChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"NukeWorkingTree":        "Dx",
		"DiscardUnstagedChanges": "Du",
		"DiscardUntrackedFiles":  "Dc",
		"DiscardStagedChanges":   "DS",
		"SoftReset":              "Ds",
		"MixedReset":             "Dm",
		"HardReset":              "Dh",
	}
	actual := map[string]string{
		"NukeWorkingTree":        cfg.Keybinding.Files.NukeWorkingTree,
		"DiscardUnstagedChanges": cfg.Keybinding.Files.DiscardUnstagedChanges,
		"DiscardUntrackedFiles":  cfg.Keybinding.Files.DiscardUntrackedFiles,
		"DiscardStagedChanges":   cfg.Keybinding.Files.DiscardStagedChanges,
		"SoftReset":              cfg.Keybinding.Files.SoftReset,
		"MixedReset":             cfg.Keybinding.Files.MixedReset,
		"HardReset":              cfg.Keybinding.Files.HardReset,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestStatusFilterChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"FilterStaged":    "<ctrl+b>s",
		"FilterUnstaged":  "<ctrl+b>u",
		"FilterTracked":   "<ctrl+b>t",
		"FilterUntracked": "<ctrl+b>T",
		"NoFilter":        "<ctrl+b>r",
	}
	actual := map[string]string{
		"FilterStaged":    cfg.Keybinding.Files.FilterStaged,
		"FilterUnstaged":  cfg.Keybinding.Files.FilterUnstaged,
		"FilterTracked":   cfg.Keybinding.Files.FilterTracked,
		"FilterUntracked": cfg.Keybinding.Files.FilterUntracked,
		"NoFilter":        cfg.Keybinding.Files.NoFilter,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestBisectChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"BisectMarkBad":       "bb",
		"BisectMarkGood":      "bg",
		"BisectSkipCurrent":   "bs",
		"BisectSkipSelected":  "bS",
		"BisectReset":         "br",
		"BisectStartMarkBad":  "bb",
		"BisectStartMarkGood": "bg",
		"BisectChooseTerms":   "bt",
	}
	actual := map[string]string{
		"BisectMarkBad":       cfg.Keybinding.Commits.BisectMarkBad,
		"BisectMarkGood":      cfg.Keybinding.Commits.BisectMarkGood,
		"BisectSkipCurrent":   cfg.Keybinding.Commits.BisectSkipCurrent,
		"BisectSkipSelected":  cfg.Keybinding.Commits.BisectSkipSelected,
		"BisectReset":         cfg.Keybinding.Commits.BisectReset,
		"BisectStartMarkBad":  cfg.Keybinding.Commits.BisectStartMarkBad,
		"BisectStartMarkGood": cfg.Keybinding.Commits.BisectStartMarkGood,
		"BisectChooseTerms":   cfg.Keybinding.Commits.BisectChooseTerms,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestGitResetChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"MixedResetToRef": "gm",
		"SoftResetToRef":  "gs",
		"HardResetToRef":  "gh",
	}
	actual := map[string]string{
		"MixedResetToRef": cfg.Keybinding.Commits.MixedResetToRef,
		"SoftResetToRef":  cfg.Keybinding.Commits.SoftResetToRef,
		"HardResetToRef":  cfg.Keybinding.Commits.HardResetToRef,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestRebaseOptionsChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"RebaseContinue": "mc",
		"RebaseAbort":    "ma",
		"RebaseSkip":     "ms",
	}
	actual := map[string]string{
		"RebaseContinue": cfg.Keybinding.Universal.RebaseContinue,
		"RebaseAbort":    cfg.Keybinding.Universal.RebaseAbort,
		"RebaseSkip":     cfg.Keybinding.Universal.RebaseSkip,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestGitFlowChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"GitFlowFinish":       "iF",
		"GitFlowStartFeature": "if",
		"GitFlowStartHotfix":  "ih",
		"GitFlowStartBugfix":  "ib",
		"GitFlowStartRelease": "ir",
	}
	actual := map[string]string{
		"GitFlowFinish":       cfg.Keybinding.Branches.GitFlowFinish,
		"GitFlowStartFeature": cfg.Keybinding.Branches.GitFlowStartFeature,
		"GitFlowStartHotfix":  cfg.Keybinding.Branches.GitFlowStartHotfix,
		"GitFlowStartBugfix":  cfg.Keybinding.Branches.GitFlowStartBugfix,
		"GitFlowStartRelease": cfg.Keybinding.Branches.GitFlowStartRelease,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestBranchDeleteChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"DeleteLocalBranch":          "dc",
		"DeleteRemoteBranch":         "dr",
		"DeleteLocalAndRemoteBranch": "db",
	}
	actual := map[string]string{
		"DeleteLocalBranch":          cfg.Keybinding.Branches.DeleteLocalBranch,
		"DeleteRemoteBranch":         cfg.Keybinding.Branches.DeleteRemoteBranch,
		"DeleteLocalAndRemoteBranch": cfg.Keybinding.Branches.DeleteLocalAndRemoteBranch,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestTagDeleteChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"DeleteLocalTag":          "dc",
		"DeleteRemoteTag":         "dr",
		"DeleteLocalAndRemoteTag": "db",
	}
	actual := map[string]string{
		"DeleteLocalTag":          cfg.Keybinding.Branches.DeleteLocalTag,
		"DeleteRemoteTag":         cfg.Keybinding.Branches.DeleteRemoteTag,
		"DeleteLocalAndRemoteTag": cfg.Keybinding.Branches.DeleteLocalAndRemoteTag,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestDiscardFilesChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"DiscardAllChanges":   "dc",
		"DiscardUnstagedFile": "du",
	}
	actual := map[string]string{
		"DiscardAllChanges":   cfg.Keybinding.Files.DiscardAllChanges,
		"DiscardUnstagedFile": cfg.Keybinding.Files.DiscardUnstagedFile,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestRebaseBranchChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"RebaseBranchSimple":      "rs",
		"RebaseBranchInteractive": "ri",
		"RebaseBranchOntoBase":    "rb",
	}
	actual := map[string]string{
		"RebaseBranchSimple":      cfg.Keybinding.Branches.RebaseBranchSimple,
		"RebaseBranchInteractive": cfg.Keybinding.Branches.RebaseBranchInteractive,
		"RebaseBranchOntoBase":    cfg.Keybinding.Branches.RebaseBranchOntoBase,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestMergeChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"MergeRegular":         "Mm",
		"MergeNonFFwd":         "Mn",
		"MergeFastForward":     "Mf",
		"MergeSquash":          "Ms",
		"MergeSquashCommitted": "MS",
	}
	actual := map[string]string{
		"MergeRegular":         cfg.Keybinding.Branches.MergeRegular,
		"MergeNonFFwd":         cfg.Keybinding.Branches.MergeNonFFwd,
		"MergeFastForward":     cfg.Keybinding.Branches.MergeFastForward,
		"MergeSquash":          cfg.Keybinding.Branches.MergeSquash,
		"MergeSquashCommitted": cfg.Keybinding.Branches.MergeSquashCommitted,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestUpstreamOptionsChordDefaults(t *testing.T) {
	cfg := GetDefaultConfig()
	expected := map[string]string{
		"ViewDivergenceFromUpstream": "ud",
		"ViewDivergenceFromBase":     "uD",
		"SetUpstream":                "us",
		"UnsetUpstream":              "uu",
		"ResetUpstreamMixed":         "ugm",
		"ResetUpstreamSoft":          "ugs",
		"ResetUpstreamHard":          "ugh",
		"RebaseUpstreamSimple":       "urs",
		"RebaseUpstreamInteractive":  "uri",
		"RebaseUpstreamOntoBase":     "urb",
	}
	actual := map[string]string{
		"ViewDivergenceFromUpstream": cfg.Keybinding.Branches.ViewDivergenceFromUpstream,
		"ViewDivergenceFromBase":     cfg.Keybinding.Branches.ViewDivergenceFromBase,
		"SetUpstream":                cfg.Keybinding.Branches.SetUpstream,
		"UnsetUpstream":              cfg.Keybinding.Branches.UnsetUpstream,
		"ResetUpstreamMixed":         cfg.Keybinding.Branches.ResetUpstreamMixed,
		"ResetUpstreamSoft":          cfg.Keybinding.Branches.ResetUpstreamSoft,
		"ResetUpstreamHard":          cfg.Keybinding.Branches.ResetUpstreamHard,
		"RebaseUpstreamSimple":       cfg.Keybinding.Branches.RebaseUpstreamSimple,
		"RebaseUpstreamInteractive":  cfg.Keybinding.Branches.RebaseUpstreamInteractive,
		"RebaseUpstreamOntoBase":     cfg.Keybinding.Branches.RebaseUpstreamOntoBase,
	}
	for name, want := range expected {
		if got := actual[name]; got != want {
			t.Errorf("default %s = %q, want %q", name, got, want)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config failed validation: %v", err)
	}
}

func TestKeybindingGroupsValidation_AliasedPrefixesAreRejected(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.Refresh = "<c-b>z"
	cfg.KeybindingGroups["files"] = map[string]KeybindingGroupConfig{
		"<c-b>":    {Name: "Filter A"},
		"<ctrl+b>": {Name: "Filter B"},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for aliased canonical prefixes, got nil")
	}
	if !strings.Contains(err.Error(), "canonicalize") {
		t.Errorf("error should mention canonicalization, got: %v", err)
	}
}

func TestKeybindingGroupsValidation_DuplicateGroupNamesAreRejected(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.Refresh = "xr"
	cfg.Keybinding.Universal.Pull = "yp"
	cfg.KeybindingGroups["files"] = map[string]KeybindingGroupConfig{
		"x": {Name: "Same name"},
		"y": {Name: "Same name"},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for duplicate group names, got nil")
	}
	if !strings.Contains(err.Error(), "same name") {
		t.Errorf("error should mention name duplication, got: %v", err)
	}
}
