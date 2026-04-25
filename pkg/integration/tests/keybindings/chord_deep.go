package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordDeep = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A three-key chord fires after all three keys are pressed",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.Pull = "abc"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))

		// Press the three-key chord — pull should fire when all three keys arrive.
		t.Views().Files().Focus().Press("abc")

		t.Views().Commits().
			Lines(
				Contains("two"),
				Contains("one"),
			)
		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
