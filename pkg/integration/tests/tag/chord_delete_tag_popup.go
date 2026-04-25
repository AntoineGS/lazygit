package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordDeleteTagPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the delete-tag chord prefix d in tags view opens a popup listing the local / remote / both delete variants.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateLightweightTag("v1.0", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			NavigateToLine(Contains("v1.0")).
			Press("d")

		t.ExpectPopup().Menu().
			Title(Equals("Delete tag")).
			ContainsLines(Contains("c").Contains("Delete local tag")).
			ContainsLines(Contains("r").Contains("Delete remote tag")).
			ContainsLines(Contains("b").Contains("Delete local and remote tag"))
	},
})
