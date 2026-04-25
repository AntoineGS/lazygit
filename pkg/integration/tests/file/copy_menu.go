package file

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// note: this is required to simulate the clipboard during CI
func expectClipboard(t *TestDriver, matcher *TextMatcher) {
	defer t.Shell().DeleteFile("clipboard")

	t.FileSystem().FileContent("clipboard", matcher)
}

var CopyMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The copy menu allows to copy name and diff of selected/all files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		// Run the test in a linked worktree so that we catch bugs where we
		// use the main repo's path instead of the current worktree's path.
		shell.EmptyCommit("initial commit")
		shell.AddWorktree("HEAD", "../linked-worktree", "mybranch")
		shell.Chdir("../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Empty panel: the chord popup opens but selecting any item is
		// a no-op. Just dismiss it.
		t.Views().Files().
			IsEmpty().
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Cancel()
			})

		t.Shell().
			CreateDir("dir").
			CreateFile("dir/1-unstaged_file", "unstaged content")

		// New untracked file: the chord menu fires diff handlers, and
		// untracked content reports as a diff (success toast).
		t.Views().Files().
			Press(keys.Universal.Refresh).
			Lines(
				Contains("dir").IsSelected(),
				Contains("unstaged_file"),
			).
			SelectNextItem()

		t.Shell().
			GitAdd("dir/1-unstaged_file").
			Commit("commit-unstaged").
			UpdateFile("dir/1-unstaged_file", "unstaged content (new)").
			CreateFileAndAdd("dir/2-staged_file", "staged content").
			Commit("commit-staged").
			UpdateFile("dir/2-staged_file", "staged content (new)").
			GitAdd("dir/2-staged_file")

		// Copy file name
		t.Views().Files().
			Press(keys.Universal.Refresh).
			Lines(
				Contains("dir"),
				Contains("unstaged_file").IsSelected(),
				Contains("staged_file"),
			).
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("File name")).
					Confirm()

				t.ExpectToast(Equals("File name copied to clipboard"))
				expectClipboard(t, Equals("1-unstaged_file"))
			})

		// Copy relative file path
		t.Views().Files().
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Relative path")).
					Confirm()

				t.ExpectToast(Equals("File path copied to clipboard"))
				expectClipboard(t, Equals("dir/1-unstaged_file"))
			})

		// Copy absolute file path
		t.Views().Files().
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Absolute path")).
					Confirm()

				t.ExpectToast(Equals("File path copied to clipboard"))
				worktreeDir, _ := os.Getwd()
				// On windows the following path would have backslashes, but we don't run integration tests on windows yet.
				expectClipboard(t, Equals(worktreeDir+"/dir/1-unstaged_file"))
			})

		// Selected path diff on a single (unstaged) file
		t.Views().Files().
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Confirm()

				t.ExpectToast(Equals("File diff copied to clipboard"))
				expectClipboard(t, Contains("+unstaged content (new)"))
			})

		// Selected path diff with staged and unstaged files (parent dir selected)
		t.Views().Files().
			SelectPreviousItem().
			Lines(
				Contains("dir").IsSelected(),
				Contains("unstaged_file"),
				Contains("staged_file"),
			).
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of selected file")).
					Confirm()

				t.ExpectToast(Equals("File diff copied to clipboard"))
				expectClipboard(t, Contains("+staged content (new)"))
			})

		// All files diff with staged files
		t.Views().Files().
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Confirm()

				t.ExpectToast(Equals("All files diff copied to clipboard"))
				expectClipboard(t, Contains("+staged content (new)"))
			})

		// All files diff with no staged files
		t.Views().Files().
			SelectNextItem().
			SelectNextItem().
			Lines(
				Contains("dir"),
				Contains("unstaged_file"),
				Contains("staged_file").IsSelected(),
			).
			Press(keys.Universal.Select).
			Press(keys.ChordPrefix.Files.CopyToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Diff of all files")).
					Confirm()

				t.ExpectToast(Equals("All files diff copied to clipboard"))
				expectClipboard(t, Contains("+staged content (new)").Contains("+unstaged content (new)"))
			})
	},
})
