package keybindings

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UniversalFetch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Universally-bound fetch runs from any view and updates the repo's fetch state",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		userCfg := cfg.GetUserConfig()
		// Use <f5> (not bound by any default view-level controller) so the
		// universal binding actually wins precedence in the commits view.
		userCfg.Keybinding.Universal.Fetch = "<f5>"
		// Disable auto-forward so master stays "behind" after fetch and we
		// can observe the ↓1 marker.
		userCfg.Git.AutoForwardBranches = "none"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		// Add a commit to the bare origin via a temporary working clone so
		// that fetching will pull in something new.
		shell.RunCommand([]string{"git", "clone", "../origin", "../origin-work"})
		shell.RunCommand([]string{"git", "-C", "../origin-work", "-c", "user.email=lazy@git.local", "-c", "user.name=lazy", "commit", "--allow-empty", "-m", "remote-only commit"})
		shell.RunCommand([]string{"git", "-C", "../origin-work", "push", "origin", "master"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Sanity: before fetch, we don't yet know about the remote-only
		// commit; status shows the repo as up to date.
		t.Views().Status().Content(Contains("✓ repo → master"))

		// Trigger universal fetch from a non-files view.
		t.Views().Commits().Focus()
		t.Views().Commits().Press("<f5>")

		// After fetch, master should be 1 commit behind its upstream.
		t.Views().Status().Content(Contains("↓1 repo → master"))
	},
})
