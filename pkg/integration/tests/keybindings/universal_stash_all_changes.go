package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UniversalStashAllChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Universally-bound stashAllChanges stashes from any view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.StashAllChanges = "S"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateFileAndAdd("a.txt", "uncommitted change")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus()
		t.Views().Branches().Press("S")
		// StashAllChanges prompts for a stash name. Confirm to accept the
		// default empty input.
		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Confirm()
		t.Views().Stash().Lines(Contains("WIP on master"))
	},
})
