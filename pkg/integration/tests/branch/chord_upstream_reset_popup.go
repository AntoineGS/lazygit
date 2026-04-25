package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordUpstreamResetPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the nested upstream-reset chord prefix ug advances the popup to a sub-popup listing the three reset variants (m/s/h).",
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

		t.GlobalPress("g")
		t.ExpectPopup().Menu().
			Title(Equals("Chord: ug …")).
			ContainsLines(Contains("m").Contains("Mixed reset to upstream")).
			ContainsLines(Contains("s").Contains("Soft reset to upstream")).
			ContainsLines(Contains("h").Contains("Hard reset to upstream"))
	},
})
