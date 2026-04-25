package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordWorkspaceResetPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the workspace-reset chord prefix D in files view opens a popup listing the seven reset variants with command-preview tooltip column.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().Press("D")

		t.ExpectPopup().Menu().
			Title(Equals("Discard / reset options")).
			ContainsLines(Contains("x").Contains("Nuke working tree")).
			ContainsLines(Contains("u").Contains("Discard unstaged changes")).
			ContainsLines(Contains("c").Contains("Discard untracked files")).
			ContainsLines(Contains("S").Contains("Discard staged changes")).
			ContainsLines(Contains("s").Contains("Soft reset")).
			ContainsLines(Contains("m").Contains("mixed reset")).
			ContainsLines(Contains("h").Contains("Hard reset")).
			// Tooltip column carries the colored command preview.
			ContainsLines(Contains("git reset --hard HEAD"))
	},
})
