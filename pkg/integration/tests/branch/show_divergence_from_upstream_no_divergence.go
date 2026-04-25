package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowDivergenceFromUpstreamNoDivergence = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show divergence from upstream when the divergence view is empty",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit1")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(Contains("master")).
			Press(keys.ChordPrefix.LocalBranches.BranchUpstreamOptions)

		t.ExpectPopup().Menu().
			Title(Equals("Branch upstream options")).
			Select(Contains("View divergence from upstream")).
			Confirm()

		t.Views().SubCommits().
			IsFocused().
			Title(Contains("Commits (master <-> origin/master)")).
			Lines(
				Contains("--- Remote ---"),
				Contains("--- Local ---"),
			)
	},
})
