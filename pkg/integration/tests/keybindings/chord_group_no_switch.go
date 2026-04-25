package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGroupNoSwitch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Group entry without switchTo advances chord and labels footer without changing view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Two leaves under <X><t>; one under <X><p>. <X> has no group entry,
		// so its footer at level-1 is just raw bindings; <X><t> has a group label.
		userCfg.Keybinding.Universal.Pull = "<X><p>"
		userCfg.Keybinding.Universal.Push = "<X><t><o>"
		userCfg.KeybindingGroups = map[string]config.KeybindingGroupConfig{
			"<X><t>": {Name: "Pull Request"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus()

		t.Views().Files().Press("X")
		// <X> has no group entry, so footer shows raw bindings under <X>:
		// p (Pull's leaf) and t (collapsed because <X><t> IS a group).
		t.Views().Options().Content(Contains("Pull Request"))
		t.Views().Files().IsFocused()

		t.Views().Files().Press("t")
		// Now pending is <X><t>; footer shows the leaf under it.
		t.Views().Options().Content(Contains("o"))
		t.Views().Files().IsFocused()

		t.GlobalPress("<esc>")
	},
})
