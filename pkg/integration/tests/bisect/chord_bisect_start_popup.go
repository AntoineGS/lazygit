package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordBisectStartPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the bisect chord prefix b in commits view (no bisect in progress) shows mark-start-bad / mark-start-good / choose-terms rows. Mid-bisect rows are hidden.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("c1")
		shell.EmptyCommit("c2")
		shell.EmptyCommit("c3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().Press("b")

		t.ExpectPopup().Menu().
			Title(Equals("Bisect options")).
			ContainsLines(Contains("b").Contains("bad")).
			ContainsLines(Contains("g").Contains("good")).
			ContainsLines(Contains("t").Contains("Choose bisect terms"))

		t.Views().Menu().Content(
			DoesNotContain("Skip current").
				DoesNotContain("Skip selected").
				DoesNotContain("Reset bisect"),
		)
	},
})
