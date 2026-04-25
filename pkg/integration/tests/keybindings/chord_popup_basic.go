package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPopupBasic = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing a chord prefix opens the popup; completing fires the binding and closes it.",
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
			Title(Equals("Chord: b …")).
			ContainsLines(Contains("p").Contains("Pull"))

		t.GlobalPress("p")

		t.Views().Status().Content(Equals("✓ repo → master"))
		t.Views().Commits().Lines(
			Contains("two"),
			Contains("one"),
		)
	},
})
