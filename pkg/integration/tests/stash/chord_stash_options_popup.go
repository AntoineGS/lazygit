package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordStashOptionsPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the stash chord prefix S opens a popup with the four stash variants; completing Ss stashes only staged changes.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-staged", "content")
		shell.CreateFileAndAdd("file-unstaged", "content")
		shell.EmptyCommit("initial commit")
		shell.UpdateFileAndAdd("file-staged", "new content")
		shell.UpdateFile("file-unstaged", "new content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("S")

		t.ExpectPopup().Menu().
			Title(Equals("Stash options")).
			ContainsLines(Contains("i").Contains("keep index")).
			ContainsLines(Contains("U").Contains("untracked")).
			// Match "Stash staged"/"Stash unstaged" to avoid the
			// "unstaged"-contains-"staged" ambiguity.
			ContainsLines(Contains("s").Contains("Stash staged")).
			ContainsLines(Contains("u").Contains("Stash unstaged"))

		t.GlobalPress("s")

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Lines(
				Contains("my stashed file"),
			)

		t.Views().Files().
			Lines(
				Equals(" M file-unstaged"),
			)
	},
})
