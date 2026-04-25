package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordIgnoreExcludePopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the ignore chord prefix i in files view opens a popup with two rows (ignore, exclude); pressing ii adds the file to .gitignore. Verifies popup stays two-column when no binding has a tooltip.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("toIgnore", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			NavigateToLine(Contains("toIgnore"))

		t.GlobalPress("i")

		t.ExpectPopup().Menu().
			Title(Equals("Ignore options")).
			ContainsLines(Contains("i").Contains(".gitignore")).
			ContainsLines(Contains("e").Contains(".git/info/exclude"))

		t.GlobalPress("i")

		t.FileSystem().FileContent(".gitignore", Equals("/toIgnore\n"))
	},
})
