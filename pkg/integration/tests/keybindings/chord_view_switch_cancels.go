package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordViewSwitchCancels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching views cancels any pending chord",
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

		// Press the chord prefix in Files...
		t.Views().Files().Focus().Press("b")

		// ...switch panels via mouse click (this bypasses the chord
		// dispatcher, going through Push/Replace which call
		// ClearPendingChord)...
		t.Views().Branches().Click(0, 0)
		t.Views().Branches().IsFocused()

		// ...and press what would have been the chord continuation.
		// The chord should already be cleared by the panel switch, so
		// pressing `p` in Branches is a fresh single-key press; nothing
		// is bound to `p` there (Pull was rebound to "bp" and Branches
		// has no single-key `p` action by default).
		t.Views().Branches().Press("p")

		// Pull never fired: status and commits unchanged.
		t.Views().Status().Content(Equals("↓1 repo → master"))
		t.Views().Commits().Lines(Contains("one"))
	},
})
