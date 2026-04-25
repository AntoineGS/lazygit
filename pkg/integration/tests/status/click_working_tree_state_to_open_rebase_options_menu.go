package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickWorkingTreeStateToOpenRebaseOptionsMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on the working tree state in the status side panel to continue the rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.Edit)

		t.Views().Status().
			Content(Contains("(rebasing) repo")).
			Click(1, 0)

		// After clicking the working-tree-state segment we should no
		// longer be mid-rebase: the click handler continues the rebase
		// to completion.
		t.Views().Status().Content(DoesNotContain("(rebasing)"))
	},
})
