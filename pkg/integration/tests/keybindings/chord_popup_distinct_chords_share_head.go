package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPopupDistinctChordsShareHead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Two chords sharing the same head key (e.g. gp and gP) both appear in the popup, even when sourced from different keybinding scopes (current-context vs. global).",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		userCfg.ChordPopupDelayMs = 0
		// Pull is global; CommitChanges is files-scoped. Both share
		// head 'g' — the popup must list both continuations.
		userCfg.Keybinding.Universal.Pull = "gP"
		userCfg.Keybinding.Files.CommitChanges = "gp"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("g")

		// 'g' is the "Reset to upstream" group in the files context, so
		// the popup uses that group's Name as title.
		t.ExpectPopup().Menu().
			Title(Equals("Reset to upstream")).
			ContainsLines(Contains("p")).
			ContainsLines(Contains("P"))
	},
})
