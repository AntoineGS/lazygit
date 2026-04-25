package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordBulkOptionsPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the bulk-options chord prefix b in submodules view opens a popup listing the four bulk variants with command-preview tooltip column.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		// No submodules required — we only verify the popup contents and
		// routing, not that the bulk action succeeds.
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().Press("b")

		t.ExpectPopup().Menu().
			Title(Equals("Bulk options")).
			ContainsLines(Contains("i").Contains("Bulk init")).
			ContainsLines(Contains("u").Contains("Bulk update")).
			ContainsLines(Contains("r").Contains("recursively")).
			ContainsLines(Contains("d").Contains("Bulk deinit")).
			// Tooltip column carries the colored command preview.
			ContainsLines(Contains("git submodule update --init"))
	},
})
