package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordUnboundCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Esc cancels the pending chord menu and returns focus to the originating view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "bp"
		userCfg.ChordPopupDelayMs = 0
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

		// Press the chord prefix → popup opens at depth 1.
		t.Views().Files().Press("b")
		t.ExpectPopup().Menu().Title(Equals("Chord: b …"))

		// Esc closes the chord menu and returns focus to Files without
		// firing any binding.
		t.GlobalPress(keys.Universal.Return)
		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
