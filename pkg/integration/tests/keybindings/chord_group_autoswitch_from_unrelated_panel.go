package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Regression: a group with switchTo whose sub-bindings live entirely in a
// view-scoped section (e.g. keybinding.branches) used to fail to enter chord
// mode when pressed from any other view, because gocui only registered the
// chord prefix for the views where a sub-binding existed. The fix routes
// switchTo prefixes through gocui's globalChordPrefixes registry so the
// auto-switch fires regardless of where the user is.
var ChordGroupAutoswitchFromUnrelatedPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "switchTo prefix with only view-scoped children still autoswitches from unrelated panels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Sub-binding lives in keybinding.branches only; nothing universal.
		// Without the fix, pressing X from the Commits view would not enter
		// chord-pending mode and the auto-switch would not fire.
		userCfg.Keybinding.Branches.RebaseBranch = "<X>r"
		userCfg.KeybindingGroups = map[string]config.KeybindingGroupConfig{
			"<X>": {Name: "Branch", SwitchTo: "localBranches"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("feature")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus()

		t.Views().Commits().Press("X")

		t.Views().Branches().IsFocused()

		// Footer should now show chord continuations from the branches scope.
		t.Views().Options().Content(Contains("r")).Content(Contains("<esc>"))

		t.GlobalPress("<esc>")
	},
})
