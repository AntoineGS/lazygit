package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordDeleteBranchPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the delete-branch chord prefix d in branches view opens a popup listing the local / remote / both delete variants.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("initial").
			NewBranch("feature").
			Checkout("master").
			CloneIntoRemote("origin").
			SetBranchUpstream("feature", "origin/feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature")).
			Press("d")

		t.ExpectPopup().Menu().
			Title(Equals("Delete branch")).
			ContainsLines(Contains("c").Contains("Delete local branch")).
			ContainsLines(Contains("r").Contains("Delete remote branch")).
			ContainsLines(Contains("b").Contains("Delete local and remote branch"))
	},
})
