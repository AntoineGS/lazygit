package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordStatusFilterPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the status filter chord prefix <c-b> in files view opens a popup listing the five filter variants. The currently active filter shows '(•)'; the others show '( )'. Selecting a filter updates the marker on the next press.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("a.txt", "staged")
		shell.CreateFile("b.txt", "unstaged")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("<c-b>")

		// Initial state: NoFilter (DisplayAll) is active. The "No filter"
		// row shows the (•) marker; the others show ( ).
		t.ExpectPopup().Menu().
			Title(Equals("Filter files")).
			ContainsLines(Contains("s").Contains("( )").Contains("Show only staged files")).
			ContainsLines(Contains("u").Contains("( )").Contains("Show only unstaged files")).
			ContainsLines(Contains("t").Contains("( )").Contains("Show only tracked files")).
			ContainsLines(Contains("T").Contains("( )").Contains("Show only untracked files")).
			ContainsLines(Contains("r").Contains("(•)").Contains("No filter"))

		// Complete the chord: press 's' to switch to staged filter.
		t.GlobalPress("s")

		// Re-open the popup; the marker should have moved to the staged row.
		t.Views().Files().Press("<c-b>")

		t.ExpectPopup().Menu().
			Title(Equals("Filter files")).
			ContainsLines(Contains("s").Contains("(•)").Contains("Show only staged files")).
			ContainsLines(Contains("u").Contains("( )").Contains("Show only unstaged files")).
			ContainsLines(Contains("t").Contains("( )").Contains("Show only tracked files")).
			ContainsLines(Contains("T").Contains("( )").Contains("Show only untracked files")).
			ContainsLines(Contains("r").Contains("( )").Contains("No filter"))
	},
})
