package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseToUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase the current branch to the selected branch upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			EmptyCommit("ensure-master").
			EmptyCommit("to-be-added"). // <- this will only exist remotely
			PushBranchAndSetUpstream("origin", "master").
			RenameCurrentBranch("master-local").
			HardReset("HEAD~1").
			NewBranchFrom("base-branch", "master-local").
			EmptyCommit("base-branch-commit").
			NewBranch("target").
			EmptyCommit("target-commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("target-commit"),
			Contains("base-branch-commit"),
			Contains("ensure-master"),
		)

		// On a branch with no upstream, just verify the popup opens.
		t.Views().Branches().
			Focus().
			Lines(
				Contains("target").IsSelected(),
				Contains("base-branch"),
				Contains("master-local"),
			).
			SelectNextItem().
			Lines(
				Contains("target"),
				Contains("base-branch").IsSelected(),
				Contains("master-local"),
			).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch upstream options")).
					Cancel()
			}).
			SelectNextItem().
			Lines(
				Contains("target"),
				Contains("base-branch"),
				Contains("master-local").IsSelected(),
			).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Branch upstream options")).
					Select(Contains("Simple rebase onto upstream")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Contains("Chord")).
					Select(Contains("Simple rebase onto upstream")).
					Confirm()
			})

		t.Views().Commits().Lines(
			Contains("target-commit"),
			Contains("base-branch-commit"),
			Contains("to-be-added"),
			Contains("ensure-master"),
		)
	},
})
