package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MergeFastForward = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Merge a branch into another using fast-forward merge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			EmptyCommit("one").
			NewBranch("branch1").
			EmptyCommit("branch1").
			Checkout("original-branch").
			NewBranchFrom("branch2", "original-branch").
			EmptyCommit("branch2").
			Checkout("original-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("branch1"),
				Contains("branch2"),
			).
			SelectNextItem().
			Press(keys.ChordPrefix.LocalBranches.Merge)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Equals("m Merge")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("branch1").IsSelected(),
				Contains("one"),
			)

		// Check that branch2 can't be merged using fast-forward
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch2")).
			Press(keys.ChordPrefix.LocalBranches.Merge)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("fast-forward")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Error")).
			Content(Contains("Cannot fast-forward 'original-branch' to 'branch2'")).
			Confirm()
	},
})
