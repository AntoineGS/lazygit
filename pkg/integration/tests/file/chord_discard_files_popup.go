package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordDiscardFilesPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the discard chord prefix d in files view opens a popup listing the discard-all and discard-unstaged variants.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.CreateFileAndAdd("staged.txt", "staged content\n")
		shell.UpdateFile("staged.txt", "staged content with unstaged change\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			NavigateToLine(Contains("staged.txt")).
			Press("d")

		t.ExpectPopup().Menu().
			Title(Equals("Discard changes")).
			ContainsLines(Contains("c").Contains("Discard")).
			ContainsLines(Contains("u").Contains("Discard unstaged"))
	},
})
