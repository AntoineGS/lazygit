package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordRebaseOptionsPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the rebase-options chord prefix m during a paused interactive rebase opens a popup with continue / abort / skip rows.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("c1").
			EmptyCommit("c2").
			EmptyCommit("c3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Pause an interactive rebase so the working tree is in REBASING state
		// and the m-chord bindings are enabled.
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("c2")).
			Press(keys.Universal.Edit)

		t.Views().Information().Content(Contains("Rebasing"))

		t.GlobalPress("m")

		t.ExpectPopup().Menu().
			Title(Equals("Rebase options")).
			ContainsLines(Contains("c").Contains("Continue rebase / merge")).
			ContainsLines(Contains("a").Contains("Abort rebase / merge")).
			ContainsLines(Contains("s").Contains("Skip current rebase commit"))
	},
})
