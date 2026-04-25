package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordEscCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Esc during a pending chord cancels it with no action fired",
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

		// Remove the 'two' commit so a pull would have something to fetch.
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Sanity: we are behind by one commit.
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))

		t.Views().Files().Focus()

		// Press the chord prefix...
		t.Views().Files().Press("b")
		// ...then Esc to cancel...
		t.Views().Files().Press("<esc>")
		// ...then what would have been the continuation. After cancel, this
		// is treated as a fresh single key press; since "bp" replaced
		// Universal.Pull, single-key "p" is unbound and nothing should fire.
		t.Views().Files().Press("p")

		// Files view is still focused and the pull never ran: status and
		// commits unchanged from setup.
		t.Views().Files().IsFocused()
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
