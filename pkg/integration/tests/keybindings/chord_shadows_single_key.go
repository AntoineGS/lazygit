package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordShadowsSingleKey = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A chord prefix shadows any single-key action; pressing the prefix alone enters chord-pending mode",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "bp"
		// Disable the popup so Files keeps focus after the prefix.
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

		t.Views().Files().Focus()

		// Pressing the chord prefix must not fire any single-key action.
		t.Views().Files().Press("b")

		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))

		t.Views().Files().Press("<esc>")
	},
})
