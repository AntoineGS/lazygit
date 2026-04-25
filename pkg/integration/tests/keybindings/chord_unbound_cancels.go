package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordUnboundCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "An unbound continuation key silently cancels the pending chord",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.Pull = "bp"
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

		t.Views().Files().Focus()
		// Press the chord prefix.
		t.Views().Files().Press("b")
		// Press an unbound continuation: should silently cancel chord.
		t.Views().Files().Press("x")

		// Still on Files; nothing fired.
		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
