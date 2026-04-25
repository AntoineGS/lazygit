package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Gitignore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that we can't ignore the .gitignore file, then ignore/exclude other files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("toExclude", "")
		shell.CreateFile("toIgnore", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ?? .gitignore"),
				Equals("  ?? toExclude"),
				Equals("  ?? toIgnore"),
			).
			SelectNextItem().
			// ensure we can't exclude the .gitignore file
			Press(keys.ChordPrefix.Files.IgnoreOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Ignore options")).
					Select(Contains("Add to .git/info/exclude")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot exclude .gitignore")).Confirm()
			}).
			// ensure we can't ignore the .gitignore file
			Press(keys.ChordPrefix.Files.IgnoreOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Ignore options")).
					Select(Contains("Add to .gitignore")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot ignore .gitignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", DoesNotContain(".gitignore"))
			}).
			SelectNextItem().
			// exclude a file
			Press(keys.ChordPrefix.Files.IgnoreOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Ignore options")).
					Select(Contains("Add to .git/info/exclude")).
					Confirm()
			}).
			Tap(func() {
				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", Contains("/toExclude"))
			}).
			SelectNextItem().
			// ignore a file
			Press(keys.ChordPrefix.Files.IgnoreOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Ignore options")).
					Select(Contains("Add to .gitignore")).
					Confirm()
			}).
			Tap(func() {
				t.FileSystem().FileContent(".gitignore", Equals("/toIgnore\n"))
				t.FileSystem().FileContent(".git/info/exclude", Contains("/toExclude"))
			})
	},
})
