package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordMenuCancelRow = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting the Cancel row in the chord menu closes it without firing any binding.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.ChordPopupDelayMs = 0
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
		t.Views().Files().Focus().Press("b")

		t.ExpectPopup().Menu().
			ContainsLines(Contains("p").Contains("Pull"))

		// Navigate to the Cancel row at the bottom and press Enter.
		t.Views().Menu().Press(keys.Universal.GotoBottom)
		t.Views().Menu().Press(keys.Universal.Confirm)

		// Menu closed; focus back on Files; no Pull fired.
		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))
	},
})
