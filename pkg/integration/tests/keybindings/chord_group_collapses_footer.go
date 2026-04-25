package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupCollapsesFooter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Footer shows a single collapsed row per group regardless of binding count",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Three bindings under <X><t>; only one collapsed row should appear.
		userCfg.Keybinding.Universal.Pull = "<X><t><o>"
		userCfg.Keybinding.Universal.Push = "<X><t><l>"
		userCfg.Keybinding.Universal.Refresh = "<X><t><r>"
		userCfg.KeybindingGroups = map[string]config.KeybindingGroupConfig{
			"<X><t>": {Name: "Pull Request"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("X")

		// One row for "Pull Request", NOT three rows for the three leaves.
		t.Views().Options().Content(Contains("Pull Request"))
		// Sanity: descriptions of the leaf bindings (Pull/Push/Refresh) should NOT appear in the
		// collapsed footer when their group is defined.
		t.Views().Options().
			Content(DoesNotContain("pull")).
			Content(DoesNotContain("push")).
			Content(DoesNotContain("refresh"))

		t.GlobalPress("<esc>")
	},
})
