package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPopupDisabled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "chordPopupDelayMs<0 disables the popup entirely; chord still dispatches.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.ChordPopupDelayMs = -1
		userCfg.Keybinding.Universal.Pull = "bp"
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

		// Press the chord prefix. With the popup disabled, no menu opens
		// and Files stays focused.
		t.Views().Files().Press("b")
		t.Views().Files().IsFocused()

		// Press the second chord key — chord still dispatches normally.
		t.Views().Files().Press("p")

		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
