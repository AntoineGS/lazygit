package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordUpstreamPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the upstream chord prefix u in branches view opens a popup listing the four simple actions plus the two nested chord prefixes (g/r).",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("initial").
			CloneIntoRemote("origin").
			SetBranchUpstream("master", "origin/master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(Contains("master")).
			Press("u")

		t.ExpectPopup().Menu().
			Title(Equals("Branch upstream options")).
			ContainsLines(Contains("d").Contains("View divergence from upstream")).
			ContainsLines(Contains("D").Contains("View divergence from base branch")).
			ContainsLines(Contains("s").Contains("Set upstream")).
			ContainsLines(Contains("u").Contains("Unset upstream")).
			// The two nested chord prefixes show as continuation rows in the popup.
			ContainsLines(Contains("g")).
			ContainsLines(Contains("r"))
	},
})
