package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordUpstreamRebasePopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the nested upstream-rebase chord prefix ur advances the popup to a sub-popup listing the three rebase variants (s/i/b).",
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
		t.Views().Branches().Focus()

		t.Views().Branches().Press("u")
		t.ExpectPopup().Menu().Title(Equals("Branch upstream options"))

		t.GlobalPress("r")
		t.ExpectPopup().Menu().
			Title(Equals("Chord: ur …")).
			ContainsLines(Contains("s").Contains("Simple rebase onto upstream")).
			ContainsLines(Contains("i").Contains("Interactive rebase onto upstream")).
			ContainsLines(Contains("b").Contains("Rebase onto base branch"))
	},
})
