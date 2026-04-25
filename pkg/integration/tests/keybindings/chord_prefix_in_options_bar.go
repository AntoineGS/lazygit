package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordPrefixInOptionsBar = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Chord-group prefixes flagged with displayOnScreen surface in the bottom options bar (replacing the menu-opener bindings' on-screen affordance).",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Files context: the 'D' chord prefix groups discard/reset
		// options; ShortName "Reset" is what shows on the bar so the
		// "Reset: D" affordance from before chording is preserved.
		t.Views().Files().Focus()
		t.Views().Options().Content(Contains("Reset: D"))

		// Local branches context: 'g' surfaces as "Reset: g" and 'u'
		// as "Upstream: u".
		t.Views().Branches().Focus()
		t.Views().Options().Content(Contains("Reset: g").Contains("Upstream: u"))
	},
})
