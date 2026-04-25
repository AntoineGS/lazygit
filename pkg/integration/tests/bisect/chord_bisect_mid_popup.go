package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordBisectMidPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the bisect chord prefix b in commits view (mid-bisect) shows mark-bad / mark-good / skip-current / reset rows. Start-mode rows are hidden.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("c1")
		shell.EmptyCommit("c2")
		shell.EmptyCommit("c3")
		shell.RunCommand([]string{"git", "bisect", "start"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Information().Content(Contains("Bisecting"))

		t.Views().Commits().Focus().Press("b")

		t.ExpectPopup().Menu().
			Title(Equals("Bisect options")).
			ContainsLines(Contains("b").Contains("bad")).
			ContainsLines(Contains("g").Contains("good")).
			ContainsLines(Contains("s").Contains("Skip current")).
			ContainsLines(Contains("r").Contains("Reset bisect"))

		t.Views().Menu().Content(DoesNotContain("(start bisect)"))
	},
})
