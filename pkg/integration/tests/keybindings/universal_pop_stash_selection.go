package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UniversalPopStashSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Universal popStash from stash view pops the SELECTED entry, not topmost",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Use <f6> to avoid colliding with the default universal Push ("P").
		userCfg.Keybinding.Universal.PopStash = "<f6>"
		userCfg.Gui.SkipStashWarning = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateFileAndAdd("first.txt", "first stashed change")
		shell.RunCommand([]string{"git", "stash", "push", "-m", "first stash"})
		shell.CreateFileAndAdd("second.txt", "second stashed change")
		shell.RunCommand([]string{"git", "stash", "push", "-m", "second stash"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().Focus().
			Lines(
				Contains("second stash"),
				Contains("first stash"),
			).
			NavigateToLine(Contains("first stash"))
		t.Views().Stash().Press("<f6>")
		// "second stash" remains; we popped the SELECTED ("first stash").
		t.Views().Stash().Lines(Contains("second stash"))
	},
})
