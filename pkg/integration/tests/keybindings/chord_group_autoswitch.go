package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupAutoswitch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing a chord prefix with switchTo activates the target context",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// Use "X" as the group prefix: uppercase X has no default binding so
		// there is no leaf/group collision with the default config.
		userCfg := cfg.GetUserConfig()
		userCfg.Keybinding.Universal.Pull = "<X><p>"
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
		// Start in commits view.
		t.Views().Commits().Focus()

		// Press the group prefix; this should auto-switch to localBranches.
		t.Views().Commits().Press("X")

		t.Views().Branches().IsFocused()

		// Footer should now show chord continuations from the new context's scope.
		t.Views().Options().Content(Contains("p")).Content(Contains("<esc>"))

		// Cleanup: cancel the chord.
		t.GlobalPress("<esc>")
	},
})
