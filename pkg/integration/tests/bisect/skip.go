package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Skip = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a git bisect and skip a few commits (selected or current)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("commit 10")).
			Press(keys.ChordPrefix.Commits.BisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect options")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()
			}).
			NavigateToLine(Contains("commit 01")).
			Press(keys.ChordPrefix.Commits.BisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect options")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
				t.Views().Information().Content(Contains("Bisecting"))
			}).
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").DoesNotContain("<--"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").DoesNotContain("<--"),
				Contains("CI commit 05").Contains("<-- current").IsSelected(),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			).
			Press(keys.ChordPrefix.Commits.BisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect options")).
					Lines(
						Contains("Mark").Contains("as bad"),
						Contains("Mark").Contains("as good"),
						Contains("Skip current commit"),
						Contains("Skip selected commit"),
						Contains("Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("Skip current commit")).Confirm()
			}).
			// Skipping the current commit selects the new current commit:
			Lines(
				Contains("CI commit 10").Contains("<-- bad"),
				Contains("CI commit 09").DoesNotContain("<--"),
				Contains("CI commit 08").DoesNotContain("<--"),
				Contains("CI commit 07").DoesNotContain("<--"),
				Contains("CI commit 06").Contains("<-- current").IsSelected(),
				Contains("CI commit 05").Contains("<-- skipped"),
				Contains("CI commit 04").DoesNotContain("<--"),
				Contains("CI commit 03").DoesNotContain("<--"),
				Contains("CI commit 02").DoesNotContain("<--"),
				Contains("CI commit 01").Contains("<-- good"),
			).
			NavigateToLine(Contains("commit 07")).
			Press(keys.ChordPrefix.Commits.BisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Bisect options")).
					Lines(
						Contains("Mark").Contains("as bad"),
						Contains("Mark").Contains("as good"),
						Contains("Skip current commit"),
						Contains("Skip selected commit"),
						Contains("Reset bisect"),
						Contains("Cancel"),
					).
					Select(Contains("Skip selected commit")).Confirm()
			}).
			// Skipping a selected, non-current commit keeps the selection
			// there:
			SelectedLine(Contains("CI commit 07").Contains("<-- skipped"))
	},
})
