package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPopupNestedRefresh = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Three-key chord: popup title and items refresh in place as the prefix grows.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.ChordPopupDelayMs = 0
		userCfg.Keybinding.Universal.Refresh = "rrr"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus()

		t.Views().Files().Press("r")
		t.ExpectPopup().Menu().Title(Equals("Chord: r …"))

		t.GlobalPress("r")
		t.ExpectPopup().Menu().Title(Equals("Chord: rr …"))

		t.GlobalPress("r")
		t.Views().Files().IsFocused()
	},
})
