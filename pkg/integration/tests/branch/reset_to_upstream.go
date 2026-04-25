package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetToUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset the current branch to the selected branch upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			NewBranch("hard-branch").
			EmptyCommit("hard commit").
			PushBranchAndSetUpstream("origin", "hard-branch").
			NewBranch("soft-branch").
			EmptyCommit("soft commit").
			PushBranchAndSetUpstream("origin", "soft-branch").
			RenameCurrentBranch("soft-branch-local").
			NewBranch("base").
			EmptyCommit("base-branch commit").
			CreateFile("file-1", "content").
			GitAdd("file-1").
			Commit("commit with file").
			CreateFile("file-2", "content").
			GitAdd("file-2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("base").IsSelected(),
				Contains("soft-branch-local"),
				Contains("hard-branch"),
			).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions).
			Tap(func() {
				// Branch has no upstream — just dismiss the popup.
				t.ExpectPopup().Menu().
					Title(Equals("Branch upstream options")).
					Cancel()
			}).
			SelectNextItem().
			Lines(
				Contains("base"),
				Contains("soft-branch-local").IsSelected(),
				Contains("hard-branch"),
			).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch upstream options")).
					Select(Contains("Mixed reset to upstream")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Contains("Chord")).
					Select(Contains("Soft reset to upstream")).
					Confirm()
			})

		t.Views().Commits().Lines(
			Contains("soft commit"),
			Contains("hard commit"),
		)
		t.Views().Files().Lines(
			Equals("▼ /"),
			Equals("  A  file-1"),
			Equals("  A  file-2"),
		)

		// hard reset
		t.Views().Branches().
			Focus().
			Lines(
				Contains("base"),
				Contains("soft-branch-local").IsSelected(),
				Contains("hard-branch"),
			).
			NavigateToLine(Contains("hard-branch")).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch upstream options")).
					Select(Contains("Mixed reset to upstream")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Contains("Chord")).
					Select(Contains("Hard reset to upstream")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Hard reset")).
					Content(Contains("Are you sure you want to do a hard reset?")).
					Confirm()
			})
		t.Views().Commits().Lines(Contains("hard commit"))
		t.Views().Files().IsEmpty()
	},
})
