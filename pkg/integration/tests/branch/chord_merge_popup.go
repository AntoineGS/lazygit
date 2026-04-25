package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordMergePopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the merge chord prefix M in branches view opens a popup listing the regular / non-ff / ff-only / squash variants.",
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
			Press("M")

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			ContainsLines(Contains("m").Contains("Merge")).
			ContainsLines(Contains("n").Contains("Regular merge (with merge commit)")).
			ContainsLines(Contains("f").Contains("Regular merge (fast-forward)")).
			ContainsLines(Contains("s").Contains("Squash merge (uncommitted)")).
			ContainsLines(Contains("S").Contains("Squash merge (committed)"))
	},
})
