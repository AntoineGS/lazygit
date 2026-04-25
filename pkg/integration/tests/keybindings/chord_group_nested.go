package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupNested = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Nested groups: X switches to localBranches, t labels Pull Request, o fires action",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "<X><t><o>"
		userCfg.KeybindingGroups = map[string]config.KeybindingGroupConfig{
			"<X>":    {Name: "Branch", SwitchTo: "localBranches"},
			"<X><t>": {Name: "Pull Request"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(Contains("one"))
		t.Views().Status().Content(Equals("↓1 repo → master"))

		// Start in files; the X prefix should switch us to localBranches.
		t.Views().Files().Focus().Press("X")
		t.Views().Branches().IsFocused()
		t.Views().Options().Content(Contains("Pull Request"))

		t.Views().Branches().Press("t")
		t.Views().Branches().IsFocused()
		t.Views().Options().Content(Contains("o"))

		t.Views().Branches().Press("o")
		// Pull fires.
		t.Views().Commits().Lines(Contains("two"), Contains("one"))
		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
