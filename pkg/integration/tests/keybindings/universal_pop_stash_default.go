package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UniversalPopStashDefault = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Universal popStash from non-stash view pops the topmost stash entry",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Use <f6> to avoid colliding with the default universal Push ("P").
		userCfg.Keybinding.Universal.PopStash = "<f6>"
		userCfg.Gui.SkipStashWarning = true // skip confirmation popup
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateFileAndAdd("first.txt", "first stashed change")
		shell.RunCommand([]string{"git", "stash", "push", "-m", "first stash"})
		shell.CreateFileAndAdd("second.txt", "second stashed change")
		shell.RunCommand([]string{"git", "stash", "push", "-m", "second stash"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().Lines(
			Contains("second stash"),
			Contains("first stash"),
		)
		t.Views().Files().Focus()
		t.Views().Files().Press("<f6>")
		// One stash remaining (the older "first stash"); topmost was popped.
		t.Views().Stash().Lines(Contains("first stash"))
	},
})
