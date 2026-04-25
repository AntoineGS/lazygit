package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordEscCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Esc during a pending chord cancels it with no action fired",
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
		t.Views().Commits().Lines(Contains("one"))

		t.Views().Files().Focus()

		t.Views().Files().Press("b")
		t.ExpectPopup().Menu().Title(Equals("Chord: b …"))

		t.GlobalPress("<esc>")
		t.Views().Files().IsFocused()

		// After cancel, "p" is a fresh single-key press; "bp" replaced
		// Universal.Pull so single-key "p" is unbound.
		t.Views().Files().Press("p")

		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
