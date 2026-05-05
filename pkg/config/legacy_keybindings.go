package config

import "strings"

// legacyChordHeadAlias declares one deprecated single-key keybinding
// field that, when set by the user, should be applied as the chord-HEAD
// prefix for a group of replacement bindings.
//
// Each alias is keyed by the *current default* chord prefix (oldPrefix).
// At config-load time, if the user explicitly set the alias field to a
// different value, every target binding's prefix is rewritten and the
// matching KeybindingGroups entry is relocated. The alias field is then
// cleared so it doesn't show up in cheatsheets or validation.
type legacyChordHeadAlias struct {
	// alias points at the deprecated yaml field on UserConfig.
	alias *string
	// oldPrefix is what the alias defaulted to before the chord refactor.
	// Target bindings whose value starts with this prefix are rewritten.
	oldPrefix string
	// targets are pointers to the new individual chord bindings whose
	// HEAD prefix (== oldPrefix) is rewritten to *alias.
	targets []*string
	// groupContexts lists the KeybindingGroups context names that hold
	// a "Name" entry keyed by oldPrefix. The entry is moved to be keyed
	// by *alias.
	groupContexts []string
}

func (cfg *UserConfig) legacyChordAliases() []legacyChordHeadAlias {
	files := &cfg.Keybinding.Files
	branches := &cfg.Keybinding.Branches
	commits := &cfg.Keybinding.Commits
	universal := &cfg.Keybinding.Universal
	submodules := &cfg.Keybinding.Submodules

	return []legacyChordHeadAlias{
		{
			alias:     &files.IgnoreFile,
			oldPrefix: "i",
			targets: []*string{
				&files.Ignore,
				&files.Exclude,
			},
			groupContexts: []string{"files"},
		},
		{
			alias:     &files.ViewStashOptions,
			oldPrefix: "S",
			targets: []*string{
				&files.StashAllChangesKeepIndex,
				&files.StashIncludeUntrackedChanges,
				&files.StashStagedChanges,
				&files.StashUnstagedChanges,
			},
			groupContexts: []string{"files"},
		},
		{
			alias:     &files.ViewResetOptions,
			oldPrefix: "D",
			targets: []*string{
				&files.NukeWorkingTree,
				&files.DiscardUnstagedChanges,
				&files.DiscardUntrackedFiles,
				&files.DiscardStagedChanges,
				&files.SoftReset,
				&files.MixedReset,
				&files.HardReset,
			},
			groupContexts: []string{"files"},
		},
		{
			alias:     &files.OpenStatusFilter,
			oldPrefix: "<ctrl+b>",
			targets: []*string{
				&files.FilterStaged,
				&files.FilterUnstaged,
				&files.FilterTracked,
				&files.FilterUntracked,
				&files.NoFilter,
			},
			groupContexts: []string{"files"},
		},
		{
			alias:     &files.CopyFileInfoToClipboard,
			oldPrefix: "y",
			targets: []*string{
				&files.CopyFileName,
				&files.CopyRelativeFilePath,
				&files.CopyAbsoluteFilePath,
				&files.CopyFileDiff,
				&files.CopyAllFilesDiff,
			},
			groupContexts: []string{"files"},
		},
		{
			alias:     &branches.RebaseBranch,
			oldPrefix: "r",
			targets: []*string{
				&branches.RebaseBranchSimple,
				&branches.RebaseBranchInteractive,
				&branches.RebaseBranchOntoBase,
			},
			groupContexts: []string{"localBranches", "remoteBranches"},
		},
		{
			alias:     &branches.MergeIntoCurrentBranch,
			oldPrefix: "M",
			targets: []*string{
				&branches.MergeRegular,
				&branches.MergeNonFFwd,
				&branches.MergeFastForward,
				&branches.MergeSquash,
				&branches.MergeSquashCommitted,
			},
			groupContexts: []string{"localBranches", "remoteBranches"},
		},
		{
			alias:     &branches.ViewGitFlowOptions,
			oldPrefix: "i",
			targets: []*string{
				&branches.GitFlowFinish,
				&branches.GitFlowStartFeature,
				&branches.GitFlowStartHotfix,
				&branches.GitFlowStartBugfix,
				&branches.GitFlowStartRelease,
			},
			groupContexts: []string{"localBranches"},
		},
		{
			alias:     &universal.CreateRebaseOptionsMenu,
			oldPrefix: "m",
			targets: []*string{
				&universal.RebaseContinue,
				&universal.RebaseAbort,
				&universal.RebaseSkip,
			},
			groupContexts: []string{"global"},
		},
		{
			alias:     &commits.ViewBisectOptions,
			oldPrefix: "b",
			targets: []*string{
				&commits.BisectMarkBad,
				&commits.BisectMarkGood,
				&commits.BisectSkipCurrent,
				&commits.BisectSkipSelected,
				&commits.BisectReset,
				&commits.BisectStartMarkBad,
				&commits.BisectStartMarkGood,
				&commits.BisectChooseTerms,
			},
			groupContexts: []string{"commits"},
		},
		{
			alias:     &commits.ViewResetOptions,
			oldPrefix: "g",
			targets: []*string{
				&commits.MixedResetToRef,
				&commits.SoftResetToRef,
				&commits.HardResetToRef,
			},
			groupContexts: []string{"commits", "reflogCommits", "subCommits"},
		},
		{
			alias:     &submodules.BulkMenu,
			oldPrefix: "b",
			targets: []*string{
				&submodules.BulkInit,
				&submodules.BulkUpdate,
				&submodules.BulkUpdateRecursive,
				&submodules.BulkDeinit,
			},
			groupContexts: []string{"submodules"},
		},
	}
}

// migrateLegacyKeybindings rewrites chord-HEAD prefixes in target
// fields and KeybindingGroups when an alias field has been set to a
// value that differs from its default (== oldPrefix).
//
// The alias field itself stays at its (possibly remapped) value so
// integration tests and cheatsheet rendering can read the chord-HEAD
// prefix off the legacy field name.
//
// Must run after yaml decode and before ResolveChordPrefixes.
func (cfg *UserConfig) migrateLegacyKeybindings() {
	for _, a := range cfg.legacyChordAliases() {
		newPrefix := *a.alias
		if newPrefix == "" || newPrefix == a.oldPrefix {
			continue
		}

		for _, target := range a.targets {
			if strings.HasPrefix(*target, a.oldPrefix) {
				*target = newPrefix + (*target)[len(a.oldPrefix):]
			}
		}
		for _, ctx := range a.groupContexts {
			ctxGroups, ok := cfg.KeybindingGroups[ctx]
			if !ok {
				continue
			}
			entry, ok := ctxGroups[a.oldPrefix]
			if !ok {
				continue
			}
			delete(ctxGroups, a.oldPrefix)
			ctxGroups[newPrefix] = entry
		}
	}
}
