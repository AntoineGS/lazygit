package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChordGitFlowPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the git-flow chord prefix i in branches view opens a popup listing finish + four start variants. The finish row shows the live branch name.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().ChordPopupDelayMs = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.NewBranch("feature/example")
		// Mark git-flow as initialized via raw config so we don't depend on the
		// `git flow init` interactive command being available in CI.
		shell.RunCommand([]string{"git", "config", "gitflow.branch.master", "main"})
		shell.RunCommand([]string{"git", "config", "gitflow.branch.develop", "develop"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.feature", "feature/"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.bugfix", "bugfix/"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.release", "release/"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.hotfix", "hotfix/"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.support", "support/"})
		shell.RunCommand([]string{"git", "config", "gitflow.prefix.versiontag", ""})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus().Press("i")

		t.ExpectPopup().Menu().
			Title(Equals("Git flow options")).
			ContainsLines(Contains("F").Contains("finish branch 'feature/example'")).
			ContainsLines(Contains("f").Contains("Start git-flow feature")).
			ContainsLines(Contains("h").Contains("Start git-flow hotfix")).
			ContainsLines(Contains("b").Contains("Start git-flow bugfix")).
			ContainsLines(Contains("r").Contains("Start git-flow release"))
	},
})
