package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OngoingOperationsPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open the ongoing-operations popup when nothing is running",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("init")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.ShowOngoingOperations)

		t.ExpectPopup().Menu().
			Title(Equals("Ongoing operations")).
			Tap(func() {
				// The prompt text is rendered as non-model lines above the menu
				// items. It appears in the menu view's buffer content.
				t.Views().Menu().Content(Contains("No operations are currently running."))
			}).
			Cancel()

		// After closing the popup, focus should return to the Files view.
		t.Views().Files().IsFocused()
	},
})
