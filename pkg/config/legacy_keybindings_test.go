package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMigrateLegacyKeybindings_noAliasSet_isNoop(t *testing.T) {
	cfg := GetDefaultConfig()

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "ii", cfg.Keybinding.Files.Ignore)
	assert.Equal(t, "ie", cfg.Keybinding.Files.Exclude)
	assert.Equal(t, "Si", cfg.Keybinding.Files.StashAllChangesKeepIndex)
	assert.Equal(t, "rs", cfg.Keybinding.Branches.RebaseBranchSimple)
	assert.Equal(t, "Mm", cfg.Keybinding.Branches.MergeRegular)
	assert.Equal(t, "mc", cfg.Keybinding.Universal.RebaseContinue)
	assert.Equal(t, "bb", cfg.Keybinding.Commits.BisectMarkBad)
	assert.Equal(t, "bi", cfg.Keybinding.Submodules.BulkInit)
	assert.Contains(t, cfg.KeybindingGroups["files"], "i")
	assert.Contains(t, cfg.KeybindingGroups["localBranches"], "M")
	assert.Contains(t, cfg.KeybindingGroups["global"], "m")
}

func TestMigrateLegacyKeybindings_aliasSameAsDefault_isNoop(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.IgnoreFile = "i" // explicit, matches default

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "i", cfg.Keybinding.Files.IgnoreFile)
	assert.Equal(t, "ii", cfg.Keybinding.Files.Ignore)
	assert.Equal(t, "ie", cfg.Keybinding.Files.Exclude)
}

func TestMigrateLegacyKeybindings_ignoreFileAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.IgnoreFile = "x"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "x", cfg.Keybinding.Files.IgnoreFile, "alias should retain user's value")
	assert.Equal(t, "xi", cfg.Keybinding.Files.Ignore)
	assert.Equal(t, "xe", cfg.Keybinding.Files.Exclude)
	assert.Contains(t, cfg.KeybindingGroups["files"], "x")
	assert.NotContains(t, cfg.KeybindingGroups["files"], "i")
	assert.Equal(t, "Ignore or exclude file", cfg.KeybindingGroups["files"]["x"].Name)
}

func TestMigrateLegacyKeybindings_defaultsKeepLegacyFieldsAtOldPrefix(t *testing.T) {
	cfg := GetDefaultConfig()

	// Defaults match old single-key bindings so integration tests can
	// reference the legacy field names directly.
	assert.Equal(t, "i", cfg.Keybinding.Files.IgnoreFile)
	assert.Equal(t, "S", cfg.Keybinding.Files.ViewStashOptions)
	assert.Equal(t, "D", cfg.Keybinding.Files.ViewResetOptions)
	assert.Equal(t, "<ctrl+b>", cfg.Keybinding.Files.OpenStatusFilter)
	assert.Equal(t, "y", cfg.Keybinding.Files.CopyFileInfoToClipboard)
	assert.Equal(t, "r", cfg.Keybinding.Branches.RebaseBranch)
	assert.Equal(t, "M", cfg.Keybinding.Branches.MergeIntoCurrentBranch)
	assert.Equal(t, "i", cfg.Keybinding.Branches.ViewGitFlowOptions)
	assert.Equal(t, "m", cfg.Keybinding.Universal.CreateRebaseOptionsMenu)
	assert.Equal(t, "b", cfg.Keybinding.Commits.ViewBisectOptions)
	assert.Equal(t, "g", cfg.Keybinding.Commits.ViewResetOptions)
	assert.Equal(t, "b", cfg.Keybinding.Submodules.BulkMenu)
}

func TestMigrateLegacyKeybindings_viewStashOptionsAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.ViewStashOptions = "z"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "zi", cfg.Keybinding.Files.StashAllChangesKeepIndex)
	assert.Equal(t, "zU", cfg.Keybinding.Files.StashIncludeUntrackedChanges)
	assert.Equal(t, "zs", cfg.Keybinding.Files.StashStagedChanges)
	assert.Equal(t, "zu", cfg.Keybinding.Files.StashUnstagedChanges)
	assert.Contains(t, cfg.KeybindingGroups["files"], "z")
	assert.NotContains(t, cfg.KeybindingGroups["files"], "S")
}

func TestMigrateLegacyKeybindings_viewResetOptionsFiles(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.ViewResetOptions = "Z"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Zx", cfg.Keybinding.Files.NukeWorkingTree)
	assert.Equal(t, "Zu", cfg.Keybinding.Files.DiscardUnstagedChanges)
	assert.Equal(t, "Zc", cfg.Keybinding.Files.DiscardUntrackedFiles)
	assert.Equal(t, "ZS", cfg.Keybinding.Files.DiscardStagedChanges)
	assert.Equal(t, "Zs", cfg.Keybinding.Files.SoftReset)
	assert.Equal(t, "Zm", cfg.Keybinding.Files.MixedReset)
	assert.Equal(t, "Zh", cfg.Keybinding.Files.HardReset)
	assert.Contains(t, cfg.KeybindingGroups["files"], "Z")
	assert.NotContains(t, cfg.KeybindingGroups["files"], "D")
}

func TestMigrateLegacyKeybindings_openStatusFilterMultiCharPrefix(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.OpenStatusFilter = "<ctrl+f>"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "<ctrl+f>s", cfg.Keybinding.Files.FilterStaged)
	assert.Equal(t, "<ctrl+f>u", cfg.Keybinding.Files.FilterUnstaged)
	assert.Equal(t, "<ctrl+f>t", cfg.Keybinding.Files.FilterTracked)
	assert.Equal(t, "<ctrl+f>T", cfg.Keybinding.Files.FilterUntracked)
	assert.Equal(t, "<ctrl+f>r", cfg.Keybinding.Files.NoFilter)
	assert.Contains(t, cfg.KeybindingGroups["files"], "<ctrl+f>")
	assert.NotContains(t, cfg.KeybindingGroups["files"], "<ctrl+b>")
}

func TestMigrateLegacyKeybindings_copyFileInfoToClipboardAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.CopyFileInfoToClipboard = "Y"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Yn", cfg.Keybinding.Files.CopyFileName)
	assert.Equal(t, "Yp", cfg.Keybinding.Files.CopyRelativeFilePath)
}

func TestMigrateLegacyKeybindings_rebaseBranchAliasUpdatesBothContexts(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Branches.RebaseBranch = "R"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Rs", cfg.Keybinding.Branches.RebaseBranchSimple)
	assert.Equal(t, "Ri", cfg.Keybinding.Branches.RebaseBranchInteractive)
	assert.Equal(t, "Rb", cfg.Keybinding.Branches.RebaseBranchOntoBase)
	assert.Contains(t, cfg.KeybindingGroups["localBranches"], "R")
	assert.Contains(t, cfg.KeybindingGroups["remoteBranches"], "R")
	assert.NotContains(t, cfg.KeybindingGroups["localBranches"], "r")
	assert.NotContains(t, cfg.KeybindingGroups["remoteBranches"], "r")
}

func TestMigrateLegacyKeybindings_mergeIntoCurrentBranchAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Branches.MergeIntoCurrentBranch = "Z"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Zm", cfg.Keybinding.Branches.MergeRegular)
	assert.Equal(t, "Zn", cfg.Keybinding.Branches.MergeNonFFwd)
	assert.Equal(t, "Zf", cfg.Keybinding.Branches.MergeFastForward)
	assert.Equal(t, "Zs", cfg.Keybinding.Branches.MergeSquash)
	assert.Equal(t, "ZS", cfg.Keybinding.Branches.MergeSquashCommitted)
	assert.Contains(t, cfg.KeybindingGroups["localBranches"], "Z")
	assert.Contains(t, cfg.KeybindingGroups["remoteBranches"], "Z")
}

func TestMigrateLegacyKeybindings_viewGitFlowOptionsAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Branches.ViewGitFlowOptions = "j"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "jF", cfg.Keybinding.Branches.GitFlowFinish)
	assert.Equal(t, "jf", cfg.Keybinding.Branches.GitFlowStartFeature)
	assert.Contains(t, cfg.KeybindingGroups["localBranches"], "j")
}

func TestMigrateLegacyKeybindings_createRebaseOptionsMenuAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Universal.CreateRebaseOptionsMenu = "Q"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Qc", cfg.Keybinding.Universal.RebaseContinue)
	assert.Equal(t, "Qa", cfg.Keybinding.Universal.RebaseAbort)
	assert.Equal(t, "Qs", cfg.Keybinding.Universal.RebaseSkip)
	assert.Contains(t, cfg.KeybindingGroups["global"], "Q")
	assert.NotContains(t, cfg.KeybindingGroups["global"], "m")
}

func TestMigrateLegacyKeybindings_viewBisectOptionsAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Commits.ViewBisectOptions = "B"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Bb", cfg.Keybinding.Commits.BisectMarkBad)
	assert.Equal(t, "Bg", cfg.Keybinding.Commits.BisectMarkGood)
	assert.Equal(t, "Bt", cfg.Keybinding.Commits.BisectChooseTerms)
	assert.Contains(t, cfg.KeybindingGroups["commits"], "B")
	assert.NotContains(t, cfg.KeybindingGroups["commits"], "b")
}

func TestMigrateLegacyKeybindings_viewResetOptionsCommitsUpdatesAllReflogContexts(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Commits.ViewResetOptions = "G"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Gm", cfg.Keybinding.Commits.MixedResetToRef)
	assert.Equal(t, "Gs", cfg.Keybinding.Commits.SoftResetToRef)
	assert.Equal(t, "Gh", cfg.Keybinding.Commits.HardResetToRef)
	assert.Contains(t, cfg.KeybindingGroups["commits"], "G")
	assert.Contains(t, cfg.KeybindingGroups["reflogCommits"], "G")
	assert.Contains(t, cfg.KeybindingGroups["subCommits"], "G")
	assert.NotContains(t, cfg.KeybindingGroups["commits"], "g")
}

func TestMigrateLegacyKeybindings_bulkMenuAlias(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.Keybinding.Submodules.BulkMenu = "B"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "Bi", cfg.Keybinding.Submodules.BulkInit)
	assert.Equal(t, "Bu", cfg.Keybinding.Submodules.BulkUpdate)
	assert.Equal(t, "Br", cfg.Keybinding.Submodules.BulkUpdateRecursive)
	assert.Equal(t, "Bd", cfg.Keybinding.Submodules.BulkDeinit)
	assert.Contains(t, cfg.KeybindingGroups["submodules"], "B")
}

func TestMigrateLegacyKeybindings_userOverridesTailValue_keepsOverride(t *testing.T) {
	// User sets the legacy chord HEAD AND a custom tail. The tail
	// override is preserved as long as it still starts with the old
	// prefix (which it does, since the user is setting it via the new
	// schema).
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.IgnoreFile = "x"
	cfg.Keybinding.Files.Ignore = "iZ" // user kept default head, custom tail

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "xZ", cfg.Keybinding.Files.Ignore)
}

func TestMigrateLegacyKeybindings_targetWithUnrelatedPrefix_isLeftAlone(t *testing.T) {
	// If the user already remapped a target binding to something
	// unrelated to the old prefix, the migration shouldn't corrupt it.
	cfg := GetDefaultConfig()
	cfg.Keybinding.Files.IgnoreFile = "x"
	cfg.Keybinding.Files.Ignore = "ja" // doesn't start with "i"

	cfg.migrateLegacyKeybindings()

	assert.Equal(t, "ja", cfg.Keybinding.Files.Ignore)
	assert.Equal(t, "xe", cfg.Keybinding.Files.Exclude) // still rewritten because "ie" starts with "i"
}

func TestUnmarshalYAML_legacyKeyTriggersMigration(t *testing.T) {
	cfg := GetDefaultConfig()
	yamlData := []byte(`
keybinding:
  files:
    ignoreFile: "x"
`)
	err := yaml.Unmarshal(yamlData, cfg)
	assert.NoError(t, err)

	assert.Equal(t, "x", cfg.Keybinding.Files.IgnoreFile, "alias should retain user's value")
	assert.Equal(t, "xi", cfg.Keybinding.Files.Ignore)
	assert.Equal(t, "xe", cfg.Keybinding.Files.Exclude)
	assert.Contains(t, cfg.KeybindingGroups["files"], "x")
	assert.NotContains(t, cfg.KeybindingGroups["files"], "i")
}

func TestUnmarshalYAML_legacyAndNewSchemaCoexist(t *testing.T) {
	cfg := GetDefaultConfig()
	yamlData := []byte(`
keybinding:
  branches:
    rebaseBranch: "R"
    rebaseBranchInteractive: "rI"
`)
	err := yaml.Unmarshal(yamlData, cfg)
	assert.NoError(t, err)

	// rebaseBranchInteractive was set explicitly with old prefix "r"
	// and legacy alias was "R"; migration rewrites "rI" → "RI".
	assert.Equal(t, "RI", cfg.Keybinding.Branches.RebaseBranchInteractive)
	assert.Equal(t, "Rs", cfg.Keybinding.Branches.RebaseBranchSimple)
	assert.Equal(t, "Rb", cfg.Keybinding.Branches.RebaseBranchOntoBase)
}
