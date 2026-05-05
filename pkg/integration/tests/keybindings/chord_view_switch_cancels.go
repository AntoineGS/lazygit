package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordViewSwitchCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching views cancels any pending chord",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "bp"
		// Disable the popup so the panel-switch handler runs on click
		// (the open menu would otherwise intercept the click).
		userCfg.ChordPopupDelayMs = -1
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↓1 repo → master"))

		t.Views().Files().Focus().Press("b")

		t.Views().Branches().Click(0, 0)
		t.Views().Branches().IsFocused()

		// Chord was cleared by the panel switch; "p" is unbound here.
		t.Views().Branches().Press("p")

		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
