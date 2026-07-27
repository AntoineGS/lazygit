package git_commands

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestObtainBranch(t *testing.T) {
	type scenario struct {
		testName                 string
		input                    []string
		storeCommitDateAsRecency bool
		expectedBranch           *models.Branch
	}

	// Use a time stamp of 2 1/2 hours ago, resulting in a recency string of "2h"
	now := time.Now().Unix()
	unixTimestamp := now - int64(2.5*60*60)
	timeStamp := strconv.FormatInt(unixTimestamp, 10)

	scenarios := []scenario{
		{
			testName:                 "TrimHeads",
			input:                    []string{"", "heads/a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				CommitUnixTimestamp: unixTimestamp,
				AheadForPull:        "?",
				BehindForPull:       "?",
				AheadForPush:        "?",
				BehindForPush:       "?",
				Head:                false,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "NoUpstream",
			input:                    []string{"", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				CommitUnixTimestamp: unixTimestamp,
				AheadForPull:        "?",
				BehindForPull:       "?",
				AheadForPush:        "?",
				BehindForPush:       "?",
				Head:                false,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "IsHead",
			input:                    []string{"*", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				CommitUnixTimestamp: unixTimestamp,
				AheadForPull:        "?",
				BehindForPull:       "?",
				AheadForPush:        "?",
				BehindForPush:       "?",
				Head:                true,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "IsBehindAndAhead",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]", "[behind 2, ahead 3]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				CommitUnixTimestamp: unixTimestamp,
				AheadForPull:        "3",
				BehindForPull:       "2",
				AheadForPush:        "3",
				BehindForPush:       "2",
				Head:                false,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "RemoteBranchIsGone",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[gone]", "[gone]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				CommitUnixTimestamp: unixTimestamp,
				UpstreamGone:        true,
				AheadForPull:        "?",
				BehindForPull:       "?",
				AheadForPush:        "?",
				BehindForPush:       "?",
				Head:                false,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "WithCommitDateAsRecency",
			input:                    []string{"", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: true,
			expectedBranch: &models.Branch{
				Name:                "a_branch",
				Recency:             "2h",
				CommitUnixTimestamp: unixTimestamp,
				AheadForPull:        "?",
				BehindForPull:       "?",
				AheadForPush:        "?",
				BehindForPush:       "?",
				Head:                false,
				Subject:             "subject",
				CommitHash:          "123",
			},
		},
		{
			testName:                 "MalformedCommitDate",
			input:                    []string{"", "a_branch", "", "", "", "subject", "123", "not-a-timestamp"},
			storeCommitDateAsRecency: true,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branch := obtainBranch(s.input, s.storeCommitDateAsRecency)
			assert.EqualValues(t, s.expectedBranch, branch)
		})
	}
}

func TestObtainBranch_EpochZeroHasRecency(t *testing.T) {
	branch := obtainBranch([]string{"", "a_branch", "", "", "", "subject", "123", "0"}, true)

	assert.Zero(t, branch.CommitUnixTimestamp)
	assert.NotEmpty(t, branch.Recency)
}

type stubBranchLoaderConfig struct{}

func (stubBranchLoaderConfig) Branches(oscommands.ICmdObjBuilder) map[string]*BranchConfig {
	return nil
}

func TestBranchLoaderHierarchy_KeepsNestedHeadInItsSubtree(t *testing.T) {
	branchesOutput := strings.Join([]string{
		rawBranchLine("", "main", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1),
		rawBranchLine("", "feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2),
		rawBranchLine("*", "head-fix", "cccccccccccccccccccccccccccccccccccccccc", 3),
	}, "\n")
	revisions := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccc",
	}
	matrixOutput := "refs/heads/main\x00" + revisions[0] + "\x000 0\x000 1\x000 2\n" +
		"refs/heads/feature\x00" + revisions[1] + "\x001 0\x000 0\x000 1\n" +
		"refs/heads/head-fix\x00" + revisions[2] + "\x002 0\x001 0\x000 0\n"
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(rawBranchesArgs("-committerdate"), branchesOutput, nil).
		ExpectGitArgs(buildBranchHierarchyForEachRefArgs(revisions)[1:], matrixOutput, nil)
	loader := hierarchyBranchLoaderForLoad(t, runner, "hierarchy", func() (BranchInfo, error) {
		return BranchInfo{}, errors.New("unexpected current branch lookup")
	})

	branches, err := loader.Load(nil, nil, nil, false, func(func() error) {}, func() {})

	assert.NoError(t, err)
	assert.Equal(t, []string{"main:0", "feature:1", "head-fix:2"}, hierarchyShape(branches))
	runner.CheckForMissingCalls()
}

func TestBranchLoaderHierarchy_FlatModesPromoteHead(t *testing.T) {
	for _, test := range []struct {
		sortOrder    string
		gitSortOrder string
	}{
		{sortOrder: "date", gitSortOrder: "-committerdate"},
		{sortOrder: "recency", gitSortOrder: "-committerdate"},
		{sortOrder: "alphabetical", gitSortOrder: "refname"},
	} {
		t.Run(test.sortOrder, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(rawBranchesArgs(test.gitSortOrder), strings.Join([]string{
					rawBranchLine("", "main", "aaaa", 1),
					rawBranchLine("*", "feature", "bbbb", 2),
				}, "\n"), nil)
			loader := hierarchyBranchLoaderForLoad(t, runner, test.sortOrder, func() (BranchInfo, error) {
				return BranchInfo{}, errors.New("unexpected current branch lookup")
			})

			branches, err := loader.Load(nil, nil, nil, false, func(func() error) {}, func() {})

			assert.NoError(t, err)
			assert.Equal(t, "feature", branches[0].Name)
			runner.CheckForMissingCalls()
		})
	}
}

func TestBranchLoaderHierarchy_ExcludesSyntheticHeadFromMatrixAtoms(t *testing.T) {
	for _, test := range []struct {
		name string
		info BranchInfo
	}{
		{name: "detached", info: BranchInfo{RefName: "detached", DisplayName: "(HEAD detached)", DetachedHead: true}},
		{name: "unborn", info: BranchInfo{RefName: "unborn"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			revisions := []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(rawBranchesArgs("-committerdate"), strings.Join([]string{
					rawBranchLine("", "main", revisions[0], 1),
					rawBranchLine("", "feature", revisions[1], 2),
				}, "\n"), nil).
				ExpectGitArgs(buildBranchHierarchyForEachRefArgs(revisions)[1:], "refs/heads/main\x00"+revisions[0]+"\x000 0\x000 1\nrefs/heads/feature\x00"+revisions[1]+"\x001 0\x000 0\n", nil)
			loader := hierarchyBranchLoaderForLoad(t, runner, "hierarchy", func() (BranchInfo, error) {
				return test.info, nil
			})

			branches, err := loader.Load(nil, nil, nil, false, func(func() error) {}, func() {})

			assert.NoError(t, err)
			assert.Equal(t, test.info.RefName, branches[0].Name)
			assert.Zero(t, branches[0].HierarchyDepth)
			runner.CheckForMissingCalls()
		})
	}
}

func TestApplyBranchHierarchy_FallsBackToFlatTimestampOrder(t *testing.T) {
	revisions := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccc",
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(buildBranchHierarchyForEachRefArgs(revisions)[1:], "", errors.New("matrix failed"))
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})
	branches := []*models.Branch{
		hierarchyBranch("main", revisions[0], 1, false),
		hierarchyBranch("feature", revisions[1], 2, false),
		hierarchyBranch("head-fix", revisions[2], 3, true),
	}
	for _, branch := range branches {
		branch.HierarchyDepth = 4
	}

	actual := loader.applyBranchHierarchy(branches)

	assert.Equal(t, []string{"head-fix:0", "feature:0", "main:0"}, hierarchyShape(actual))
	runner.CheckForMissingCalls()
}

func TestApplyBranchHierarchy_FallsBackOnStructurallyIncompleteModernOutput(t *testing.T) {
	revisions := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccc",
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00"+revisions[0]+"\x000 0\x000 1\x000 2\n"+
				"refs/heads/feature\x00"+revisions[1]+"\x001 0\x000 0\x000 1\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})
	branches := []*models.Branch{
		hierarchyBranch("main", revisions[0], 1, false),
		hierarchyBranch("feature", revisions[1], 2, false),
		hierarchyBranch("head-fix", revisions[2], 3, true),
	}

	actual := loader.applyBranchHierarchy(branches)

	assert.Equal(t, []string{"head-fix:0", "feature:0", "main:0"}, hierarchyShape(actual))
	runner.CheckForMissingCalls()
}

func TestApplyBranchHierarchy_DoesNotLoseBranchesWhenRefsSwap(t *testing.T) {
	revisions := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00"+revisions[1]+"\x000 0\x001 0\n"+
				"refs/heads/feature\x00"+revisions[0]+"\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})
	branches := []*models.Branch{
		hierarchyBranch("main", revisions[0], 1, false),
		hierarchyBranch("feature", revisions[1], 2, true),
	}

	actual := loader.applyBranchHierarchy(branches)

	assert.Equal(t, []string{"feature:0", "main:0"}, hierarchyShape(actual))
	runner.CheckForMissingCalls()
}

func TestApplyBranchHierarchy_FallsBackOnMalformedConcurrentLegacyRef(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"},
			"refs/heads/bad ref\x00dddd\nrefs/heads/main\x00aaaa\n",
			nil,
		).
		ExpectGitArgs([]string{"check-ref-format", "refs/heads/bad ref"}, "", errors.New("invalid ref"))
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaa", 1, false),
		hierarchyBranch("feature", "bbbb", 2, false),
		hierarchyBranch("head-fix", "cccc", 3, true),
	}

	actual := loader.applyBranchHierarchy(branches)

	assert.Equal(t, []string{"head-fix:0", "feature:0", "main:0"}, hierarchyShape(actual))
	runner.CheckForMissingCalls()
}

func hierarchyBranchLoaderForLoad(
	t *testing.T,
	runner *oscommands.FakeCmdObjRunner,
	sortOrder string,
	getCurrentBranchInfo func() (BranchInfo, error),
) *BranchLoader {
	t.Helper()
	userConfig := config.GetDefaultConfig()
	userConfig.Git.LocalBranchSortOrder = sortOrder
	gitCommon := buildGitCommon(commonDeps{
		runner: runner, userConfig: userConfig, gitVersion: &GitVersion{2, 41, 0, ""},
	})
	return &BranchLoader{
		Common:               gitCommon.Common,
		GitCommon:            gitCommon,
		cmd:                  gitCommon.cmd,
		getCurrentBranchInfo: getCurrentBranchInfo,
		config:               stubBranchLoaderConfig{},
	}
}

func rawBranchLine(head, name, hash string, timestamp int64) string {
	return strings.Join([]string{head, "heads/" + name, "", "", "", "subject", hash, strconv.FormatInt(timestamp, 10)}, "\x00")
}

func rawBranchesArgs(sortOrder string) []string {
	return []string{
		"for-each-ref",
		"--sort=" + sortOrder,
		"--format=%(HEAD)%00%(refname:short)%00%(upstream:short)%00%(upstream:track)%00%(push:track)%00%(subject)%00%(objectname)%00%(committerdate:unix)",
		"refs/heads",
	}
}

func TestGetBehindBaseBranchValuesForAllBranches_FastPath(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master", "refs/remotes/origin/develop"}

	// Two branches: feat-x has clear divergence from develop; main matches master exactly.
	branches := []*models.Branch{
		{Name: "feat-x"},
		{Name: "main"},
	}

	expectedFormat := "%(refname)%00%(ahead-behind:refs/heads/master)%00%(ahead-behind:refs/remotes/origin/develop)"
	output := "refs/heads/feat-x\x0055 0\x005 5\n" + // picks develop (ahead=5 < 55), behind=5
		"refs/heads/main\x000 0\x000 0\n" // picks master (first, tie), behind=0

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--format=" + expectedFormat, "refs/heads"}, output, nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 41, 0, ""},
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	rendered := false
	err := loader.GetBehindBaseBranchValuesForAllBranches(branches, mainBranches, func() { rendered = true })
	assert.NoError(t, err)
	assert.True(t, rendered, "renderFunc should have been called")

	assert.Equal(t, int32(5), branches[0].BehindBaseBranch.Load(), "feat-x should be behind develop by 5")
	assert.Equal(t, int32(0), branches[1].BehindBaseBranch.Load(), "main should be behind master by 0")

	runner.CheckForMissingCalls()
}

// edge case where a failure would leave artifacts from prior load
func TestGetBehindBaseBranchValuesForAllBranches_FastPath_ClearsStaleValueWhenBranchMissingFromOutput(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master"}

	feat := &models.Branch{Name: "feat-x"}
	feat.BehindBaseBranch.Store(99) // stale value from a prior load
	ghost := &models.Branch{Name: "ghost"}
	ghost.BehindBaseBranch.Store(42) // stale value from a prior load

	expectedFormat := "%(refname)%00%(ahead-behind:refs/heads/master)"
	output := "refs/heads/feat-x\x003 5\n" // ghost is intentionally absent

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--format=" + expectedFormat, "refs/heads"}, output, nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 41, 0, ""},
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	err := loader.GetBehindBaseBranchValuesForAllBranches(
		[]*models.Branch{feat, ghost}, mainBranches, func() {})
	assert.NoError(t, err)

	assert.Equal(t, int32(5), feat.BehindBaseBranch.Load(), "feat-x should be updated to fresh value")
	assert.Equal(t, int32(0), ghost.BehindBaseBranch.Load(), "ghost should be reset to 0 since it has no fresh data")

	runner.CheckForMissingCalls()
}

func TestGetBehindBaseBranchValuesForAllBranches_FastPath_ClearsStaleValueWhenAllFieldsAreInvalid(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master"}

	feat := &models.Branch{Name: "feat-x"}
	feat.BehindBaseBranch.Store(99)

	expectedFormat := "%(refname)%00%(ahead-behind:refs/heads/master)"
	output := "refs/heads/feat-x\x00not-a-count\n"

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--format=" + expectedFormat, "refs/heads"}, output, nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 41, 0, ""},
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	err := loader.GetBehindBaseBranchValuesForAllBranches([]*models.Branch{feat}, mainBranches, func() {})
	assert.NoError(t, err)
	assert.Equal(t, int32(0), feat.BehindBaseBranch.Load())

	runner.CheckForMissingCalls()
}

func TestGetBehindBaseBranchValuesForAllBranches_LegacyPath(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master"}

	branches := []*models.Branch{
		{Name: "feat-x"},
	}

	// In legacy path: per-branch GetBaseBranch (merge-base + for-each-ref --contains)
	// then rev-list --left-right --count.
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"merge-base", "refs/heads/feat-x", "refs/heads/master"}, "abc123\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--contains", "abc123", "--format=%(refname)", "refs/heads/master"}, "refs/heads/master\n", nil).
		ExpectGitArgs([]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/master"}, "5\t7\n", nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 34, 0, ""}, // pre-2.41, forces legacy
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	rendered := false
	err := loader.GetBehindBaseBranchValuesForAllBranches(branches, mainBranches, func() { rendered = true })
	assert.NoError(t, err)
	assert.True(t, rendered)
	assert.Equal(t, int32(7), branches[0].BehindBaseBranch.Load())

	runner.CheckForMissingCalls()
}
