package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGitResetPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the git-reset chord prefix g in branches view opens a popup listing the three reset variants (mixed/soft/hard) with red command-preview tooltip column.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("current-branch")
		shell.EmptyCommit("root commit")

		shell.NewBranch("other-branch")
		shell.EmptyCommit("other-branch commit")

		shell.Checkout("current-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus().
			Lines(
				Contains("current-branch").IsSelected(),
				Contains("other-branch"),
			).
			SelectNextItem().
			Press(keys.ChordPrefix.Get("localBranches", config.ChordIDResetToRef))

		// The title is dynamically interpolated with the selected branch name.
		t.ExpectPopup().Menu().
			Title(Contains("Reset")).
			ContainsLines(Contains("m").Contains("Mixed reset")).
			ContainsLines(Contains("s").Contains("Soft reset")).
			ContainsLines(Contains("h").Contains("Hard reset")).
			// Tooltip column carries the colored command preview for the selected branch.
			ContainsLines(Contains("reset --hard other-branch"))
	},
})
