package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordRebaseBranchPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the rebase chord prefix r in branches view opens a popup listing simple / interactive / onto-base-branch variants.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("initial").
			NewBranch("feature").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature")).
			Press("r")

		t.ExpectPopup().Menu().
			Title(Equals("Rebase options")).
			ContainsLines(Contains("s").Contains("Rebase")).
			ContainsLines(Contains("i").Contains("Interactive rebase")).
			ContainsLines(Contains("b").Contains("Rebase onto base"))
	},
})
