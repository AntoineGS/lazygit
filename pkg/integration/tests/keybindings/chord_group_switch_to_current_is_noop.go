package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupSwitchToCurrentIsNoop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "switchTo target equal to current context advances chord without disruption",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "<X><p>"
		userCfg.KeybindingGroups = map[string]config.KeybindingGroupConfig{
			"<X>": {Name: "Branch", SwitchTo: "localBranches"},
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
		// Start ALREADY in localBranches.
		t.Views().Branches().Focus()

		// Press the prefix; switch is a no-op (already there) but chord advances.
		t.Views().Branches().Press("X")
		t.Views().Branches().IsFocused()
		t.Views().Options().Content(Contains("p")).Content(Contains("<esc>"))

		// Complete the chord — Pull fires.
		t.Views().Branches().Press("p")
		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
