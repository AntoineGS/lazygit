package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupCollapsesPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Popup shows a single collapsed row per group regardless of binding count",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.ChordPopupDelayMs = 0
		// Three bindings under <X><t>; only one collapsed row should appear.
		userCfg.Keybinding.Universal.Pull = "<X><t><o>"
		userCfg.Keybinding.Universal.Push = "<X><t><l>"
		userCfg.Keybinding.Universal.Refresh = "<X><t><r>"
		userCfg.KeybindingGroups = map[string]map[string]config.KeybindingGroupConfig{
			"global": {
				"<X><t>": {Name: "Pull Request"},
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("X")

		// Popup is up at depth 1; the 't' continuation collapses the three
		// leaves under "Pull Request" into a single row.
		t.ExpectPopup().Menu().
			Title(Equals("Chord: X …")).
			ContainsLines(Contains("t").Contains("Pull Request"))
	},
})
