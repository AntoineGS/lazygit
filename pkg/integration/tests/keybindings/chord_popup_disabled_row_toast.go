package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPopupDisabledRowToast = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing a key for a chord-popup row whose binding has a DisabledReason short-circuits with the standard disabled toast (matching the regular dispatch path) instead of running the handler.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Only one branch (master). Selecting it and trying to merge it
		// into itself triggers the notMergingIntoYourself disabled
		// reason on every M-prefixed merge variant.
		t.Views().Branches().
			Focus().
			Press("M")

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Merge")).
			Confirm()

		// Disabled toast shows; menu stays open (no merge happens).
		t.ExpectToast(Contains("You cannot merge a branch into itself"))
		t.Views().Menu().IsVisible()
	},
})
