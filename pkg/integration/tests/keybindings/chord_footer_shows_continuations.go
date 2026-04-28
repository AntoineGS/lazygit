package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordFooterShowsContinuations = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Options footer shows chord continuations while a prefix is pending",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Keybinding.Universal.Pull = "bp"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus()

		// Before pressing anything, the footer should show the normal
		// (non-chord) suggestions.
		t.Views().Options().Content(DoesNotContain("<esc>"))

		// Press the chord prefix.
		t.Views().Files().Press("b")

		// While the chord is pending, the options footer must show the
		// continuation key label and the cancel hint.
		t.Views().Options().Content(Contains("p"))
		t.Views().Options().Content(Contains("<esc>"))

		// Cancel chord state for cleanup.
		t.Views().Files().Press("<esc>")
	},
})
