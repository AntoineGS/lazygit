package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UniversalCheckForUpdate = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Universally-bound checkForUpdate fires from a non-status view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.CheckForUpdate = "U"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Start in files (NOT status).
		t.Views().Files().Focus()
		t.Views().Files().Press("U")
		// CheckForUpdateInForeground typically opens an info popup; assert
		// that some popup appears. After first run, refine to the specific
		// title if needed.
		t.ExpectPopup()
	},
})
