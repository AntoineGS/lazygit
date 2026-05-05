package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordMenuNavigable = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Chord menu is navigable: arrow keys move selection, Enter on highlighted row fires the binding.",
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
		t.Views().Files().Focus()

		t.Views().Files().Press("b")

		t.ExpectPopup().Menu().
			ContainsLines(Contains("p").Contains("Pull"))

		// Navigate. The point is that arrow keys move selection rather
		// than being silent-cancelled by chord dispatch.
		t.Views().Menu().Press(keys.Universal.NextItemAlt)
		t.Views().Menu().Press(keys.Universal.PrevItemAlt)

		t.Views().Menu().Press(keys.Universal.Confirm)

		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
