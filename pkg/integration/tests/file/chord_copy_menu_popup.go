package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordCopyMenuPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the copy chord prefix y in files view opens a popup with five rows; popup is two-column (Tooltip text is reachable via help mode, not inline).",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
		cfg.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		// Stage a file so canCopyFileDiff/canCopyAllFilesDiff are satisfied
		// and the ys/ya rows are present in the popup.
		shell.CreateFileAndAdd("file1", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		defer t.Shell().DeleteFile("clipboard")

		t.Views().Files().Focus().Press("y")

		t.ExpectPopup().Menu().
			Title(Equals("Copy to clipboard")).
			ContainsLines(Contains("n").Contains("File name")).
			ContainsLines(Contains("p").Contains("Relative path")).
			ContainsLines(Contains("P").Contains("Absolute path")).
			ContainsLines(Contains("s").Contains("Diff of selected file")).
			ContainsLines(Contains("a").Contains("Diff of all files"))

		t.GlobalPress("n")

		t.ExpectToast(Equals("File name copied to clipboard"))
		t.FileSystem().FileContent("clipboard", Equals("file1"))
	},
})
