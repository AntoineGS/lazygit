package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordShadowsSingleKey = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A chord prefix shadows any single-key action; pressing the prefix alone enters chord-pending mode",
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

		// Press the chord prefix alone — must NOT trigger any single-key
		// action. The spec invariant is that any chord starting with `b`
		// makes the single-key `b` unreachable.
		t.Views().Files().Press("b")

		// Confirm chord-pending state is active by checking the options
		// footer shows the continuation hint plus <esc>.
		t.Views().Options().Content(Contains("p"))
		t.Views().Options().Content(Contains("<esc>"))

		// Files view is still focused; no popup opened.
		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))

		// Cancel before exiting the test.
		t.Views().Files().Press("<esc>")
	},
})
