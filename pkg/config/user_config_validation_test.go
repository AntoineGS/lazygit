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
				// Multi-rune strings are now valid as chord sequences;
				// use an unterminated bracket to assert a rejected key.
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
	cfg.KeybindingGroups = map[string]KeybindingGroupConfig{
		"<bogus": {Name: "Bad"},
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
	cfg.KeybindingGroups = map[string]KeybindingGroupConfig{
		"<X>": {Name: ""},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "non-empty name") {
		t.Fatalf("expected validation error for empty name, got: %v", err)
	}
}

func TestKeybindingGroup_LeafCollisionRejected(t *testing.T) {
	cfg := GetDefaultConfig()
	// Bind exactly <b> as a leaf (single-key chord-equivalent).
	cfg.Keybinding.Universal.Pull = "<b>"
	// AND declare <b> as a group with at least one child binding.
	cfg.Keybinding.Universal.Push = "<b><p>"
	cfg.KeybindingGroups = map[string]KeybindingGroupConfig{
		"<b>": {Name: "Branch"},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "<b>") {
		t.Fatalf("expected leaf/group collision error for <b>, got: %v", err)
	}
}

func TestKeybindingGroup_MustHaveAtLeastOneBinding(t *testing.T) {
	cfg := GetDefaultConfig()
	// Note: no chord binding starts with <z>.
	cfg.KeybindingGroups = map[string]KeybindingGroupConfig{
		"<z>": {Name: "Empty"},
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "<z>") {
		t.Fatalf("expected error citing empty group <z>, got: %v", err)
	}
}

func TestKeybindingUniversal_PopStashAccepted(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.PopStash = "sp"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation to succeed for allowlisted popStash, got: %v", err)
	}
}

func TestKeybindingUniversal_NotEligibleRejected(t *testing.T) {
	// Temporarily simulate "popStash not yet annotated" by clearing the allowlist.
	original := universalEligibleActions
	universalEligibleActions = map[string]bool{}
	defer func() { universalEligibleActions = original }()

	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.PopStash = "sp"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for unannotated universal popStash")
	}
	if !strings.Contains(err.Error(), "PopStash") && !strings.Contains(err.Error(), "popStash") {
		t.Fatalf("error should cite the action name, got: %v", err)
	}
}
