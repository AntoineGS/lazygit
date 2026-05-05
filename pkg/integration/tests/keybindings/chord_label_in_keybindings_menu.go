package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordLabelInKeybindingsMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The keybindings menu (?) shows the full chord label, not just the head key.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Files.CommitChanges = "gp"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.OptionMenu)

		// "@" filters by item.Key. Searching "@gp" must match the
		// full chord label, not just the head "g".
		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Filter("@gp").
			ContainsLines(Contains("gp").Contains("Commit"))
	},
})
