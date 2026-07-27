package git_commands

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func hierarchyBranch(name, hash string, timestamp int64, head bool) *models.Branch {
	return &models.Branch{
		Name: name, CommitHash: hash, CommitUnixTimestamp: timestamp, Head: head,
	}
}

func relation(child string, candidates map[string]int) branchRelationMatrix {
	return branchRelationMatrix{"refs/heads/" + child: candidates}
}

func hierarchyShape(branches []*models.Branch) []string {
	return lo.Map(branches, func(branch *models.Branch, _ int) string {
		return fmt.Sprintf("%s:%d", branch.Name, branch.HierarchyDepth)
	})
}

func mergeRelations(relations ...branchRelationMatrix) branchRelationMatrix {
	result := branchRelationMatrix{}
	for _, matrix := range relations {
		for child, candidates := range matrix {
			result[child] = candidates
		}
	}
	return result
}

func TestInferBranchHierarchy(t *testing.T) {
	tests := []struct {
		name            string
		branches        []*models.Branch
		relations       branchRelationMatrix
		mainBranchNames []string
		expected        []string
	}{
		{
			name: "linear",
			branches: []*models.Branch{
				hierarchyBranch("main", "a", 1, false),
				hierarchyBranch("feature", "b", 2, false),
				hierarchyBranch("fix", "c", 3, false),
			},
			relations: mergeRelations(
				relation("feature", map[string]int{"refs/heads/main": 1}),
				relation("fix", map[string]int{"refs/heads/feature": 1}),
			),
			expected: []string{"main:0", "feature:1", "fix:2"},
		},
		{
			name: "siblings by date",
			branches: []*models.Branch{
				hierarchyBranch("main", "a", 1, false),
				hierarchyBranch("older", "b", 2, false),
				hierarchyBranch("newer", "c", 3, false),
			},
			relations: mergeRelations(
				relation("older", map[string]int{"refs/heads/main": 1}),
				relation("newer", map[string]int{"refs/heads/main": 1}),
			),
			expected: []string{"main:0", "newer:1", "older:1"},
		},
		{
			name: "unrelated roots",
			branches: []*models.Branch{
				hierarchyBranch("older-root", "a", 1, false),
				hierarchyBranch("newer-root", "b", 2, false),
			},
			expected: []string{"newer-root:0", "older-root:0"},
		},
		{
			name: "same-tip refs",
			branches: []*models.Branch{
				hierarchyBranch("alias", "a", 1, false),
				hierarchyBranch("main", "a", 2, false),
			},
			expected: []string{"main:0", "alias:0"},
		},
		{
			name: "nested HEAD subtree",
			branches: []*models.Branch{
				hierarchyBranch("newer-unrelated-root", "a", 3, false),
				hierarchyBranch("old-root", "b", 1, false),
				hierarchyBranch("child-head", "c", 2, true),
			},
			relations: relation("child-head", map[string]int{"refs/heads/old-root": 1}),
			expected:  []string{"old-root:0", "child-head:1", "newer-unrelated-root:0"},
		},
		{
			name: "detached HEAD",
			branches: []*models.Branch{
				{Name: "detached", DisplayName: "(HEAD detached)", Head: true, DetachedHead: true},
				hierarchyBranch("feature", "b", 2, false),
				hierarchyBranch("main", "a", 1, false),
			},
			relations: relation("feature", map[string]int{"refs/heads/main": 1}),
			expected:  []string{"detached:0", "main:0", "feature:1"},
		},
		{
			name:     "unborn HEAD",
			branches: []*models.Branch{{Name: "unborn", Head: true}},
			expected: []string{"unborn:0"},
		},
		{
			name: "transitive reduction",
			branches: []*models.Branch{
				hierarchyBranch("main", "a", 1, false),
				hierarchyBranch("feature", "b", 2, false),
				hierarchyBranch("fix", "c", 3, false),
			},
			relations: mergeRelations(
				relation("feature", map[string]int{"refs/heads/main": 1}),
				relation("fix", map[string]int{"refs/heads/main": 2, "refs/heads/feature": 1}),
			),
			expected: []string{"main:0", "feature:1", "fix:2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := inferBranchHierarchy(test.branches, test.relations, test.mainBranchNames)
			assert.Equal(t, test.expected, hierarchyShape(actual))
		})
	}
}

func TestInferBranchHierarchy_ParentTieBreakers(t *testing.T) {
	tests := []struct {
		name            string
		branches        []*models.Branch
		relations       branchRelationMatrix
		mainBranchNames []string
		expectedParent  string
	}{
		{
			name: "candidate with more unambiguous children",
			branches: []*models.Branch{
				hierarchyBranch("a", "a", 1, false),
				hierarchyBranch("b", "b", 1, false),
				hierarchyBranch("a-child", "c", 2, false),
				hierarchyBranch("target", "d", 3, false),
			},
			relations: mergeRelations(
				relation("a-child", map[string]int{"refs/heads/a": 1}),
				relation("target", map[string]int{"refs/heads/a": 1, "refs/heads/b": 1}),
			),
			expectedParent: "a",
		},
		{
			name: "configured main branch",
			branches: []*models.Branch{
				hierarchyBranch("a", "a", 1, false),
				hierarchyBranch("main", "b", 1, false),
				hierarchyBranch("target", "c", 2, false),
			},
			relations:       relation("target", map[string]int{"refs/heads/a": 1, "refs/heads/main": 1}),
			mainBranchNames: []string{"main"},
			expectedParent:  "main",
		},
		{
			name: "older candidate",
			branches: []*models.Branch{
				hierarchyBranch("newer", "a", 2, false),
				hierarchyBranch("older", "b", 1, false),
				hierarchyBranch("target", "c", 3, false),
			},
			relations:      relation("target", map[string]int{"refs/heads/newer": 1, "refs/heads/older": 1}),
			expectedParent: "older",
		},
		{
			name: "alphabetically smaller candidate",
			branches: []*models.Branch{
				hierarchyBranch("beta", "a", 1, false),
				hierarchyBranch("alpha", "b", 1, false),
				hierarchyBranch("target", "c", 2, false),
			},
			relations:      relation("target", map[string]int{"refs/heads/beta": 1, "refs/heads/alpha": 1}),
			expectedParent: "alpha",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ordered := inferBranchHierarchy(test.branches, test.relations, test.mainBranchNames)
			targetIndex := lo.IndexOf(hierarchyShape(ordered), "target:1")
			assert.Positive(t, targetIndex)
			assert.Equal(t, test.expectedParent, ordered[targetIndex-1].Name)
		})
	}
}

func TestInferBranchHierarchy_IsDeterministicAcrossInputPermutations(t *testing.T) {
	relations := relation("target", map[string]int{"refs/heads/beta": 1, "refs/heads/alpha": 1})
	alpha := hierarchyBranch("alpha", "a", 1, false)
	beta := hierarchyBranch("beta", "b", 1, false)
	target := hierarchyBranch("target", "c", 2, false)

	expected := []string{"alpha:0", "target:1", "beta:0"}
	assert.Equal(t, expected, hierarchyShape(inferBranchHierarchy(
		[]*models.Branch{alpha, beta, target}, relations, nil)))
	assert.Equal(t, expected, hierarchyShape(inferBranchHierarchy(
		[]*models.Branch{target, beta, alpha}, relations, nil)))
}

func TestBatchAheadBehindBaseRevisions(t *testing.T) {
	revisions := make([]string, 0, 800)
	for i := 1; i <= 400; i++ {
		revisions = append(revisions, fmt.Sprintf("%040x", i), fmt.Sprintf("%064x", i))
	}

	batches := batchAheadBehindBaseRevisions(revisions)
	actual := make([]string, 0, len(revisions))
	for _, batch := range batches {
		formatArg := buildBranchHierarchyForEachRefArgs(batch)[2]
		assert.LessOrEqual(t, len(formatArg), maxAheadBehindFormatArgBytes)
		actual = append(actual, batch...)
	}
	assert.Equal(t, revisions, actual)
}

func TestLoadBranchRelationMatrixFast_UsesSafeOIDAtoms(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main)", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, true),
	}
	format := "%(refname)%00%(objectname)%00%(ahead-behind:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)%00%(ahead-behind:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb)"
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--format=" + format, "refs/heads"},
			"refs/heads/main)\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\x001 0\n"+
				"refs/heads/feature\x00bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.NoError(t, err)
	assert.Equal(t, 1, matrix["refs/heads/feature"]["refs/heads/main)"])
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixFast_ExpandsDuplicateTips(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("parent-one", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("parent-two", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 2, false),
		hierarchyBranch("child", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 3, false),
	}
	format := "%(refname)%00%(objectname)%00%(ahead-behind:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)%00%(ahead-behind:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb)"
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--format=" + format, "refs/heads"},
			"refs/heads/parent-one\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\x000 1\n"+
				"refs/heads/parent-two\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\x000 1\n"+
				"refs/heads/child\x00bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.NoError(t, err)
	candidates := slices.Sorted(slices.Values(lo.Keys(matrix["refs/heads/child"])))
	assert.Equal(t, []string{"refs/heads/parent-one", "refs/heads/parent-two"}, candidates)
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixFast_DiscardsPartialMatrixOnError(t *testing.T) {
	revisions := make([]string, 0, 500)
	branches := make([]*models.Branch, 0, 500)
	for i := 1; i <= 500; i++ {
		revision := fmt.Sprintf("%040x", i)
		revisions = append(revisions, revision)
		branches = append(branches, hierarchyBranch(fmt.Sprintf("branch-%d", i), revision, int64(i), false))
	}
	batches := batchAheadBehindBaseRevisions(revisions)
	assert.Greater(t, len(batches), 1)
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(buildBranchHierarchyForEachRefArgs(batches[0])[1:], modernZeroOutput(branches, len(batches[0])), nil).
		ExpectGitArgs(buildBranchHierarchyForEachRefArgs(batches[1])[1:], "", errors.New("matrix failed"))
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "matrix failed")
	runner.CheckForMissingCalls()
}

func modernZeroOutput(branches []*models.Branch, cellCount int) string {
	cells := strings.Repeat("\x000 0", cellCount)
	return strings.Join(lo.Map(branches, func(branch *models.Branch, _ int) string {
		return branch.FullRefName() + "\x00" + branch.CommitHash + cells
	}), "\n")
}

func TestLoadBranchRelationMatrixFast_RejectsMissingExpectedRows(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, false),
	}
	revisions := []string{branches[0].CommitHash, branches[1].CommitHash}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\x000 1\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "branch hierarchy output missing row for refs/heads/feature")
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixFast_RejectsWrongColumnCount(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, false),
	}
	revisions := []string{branches[0].CommitHash, branches[1].CommitHash}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\n"+
				"refs/heads/feature\x00bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "malformed branch hierarchy row for refs/heads/main: expected 4 columns, got 3")
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixFast_ToleratesMalformedCellsWithoutShiftingBaseIdentity(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("parent", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, false),
		hierarchyBranch("feature", "cccccccccccccccccccccccccccccccccccccccc", 3, false),
	}
	revisions := []string{branches[0].CommitHash, branches[1].CommitHash, branches[2].CommitHash}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x000 0\x00\x000 2\n"+
				"refs/heads/parent\x00bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\x00malformed\x000 0\x000 1\n"+
				"refs/heads/feature\x00cccccccccccccccccccccccccccccccccccccccc\x00malformed\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"refs/heads/parent": 1}, matrix["refs/heads/feature"])
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixFast_RejectsMovedKnownRef(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
		hierarchyBranch("feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, false),
	}
	format := "%(refname)%00%(objectname)%00%(ahead-behind:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)%00%(ahead-behind:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb)"
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--format=" + format, "refs/heads"},
			"refs/heads/main\x00bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\x000 0\x001 0\n"+
				"refs/heads/feature\x00aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixFast(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "branch hierarchy snapshot mismatch for refs/heads/main: expected aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, got bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	runner.CheckForMissingCalls()
}

func TestLoadDirectBranchCandidatesFast_RetainsStrictAncestorsForDistanceSelection(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaa", 1, false),
		hierarchyBranch("feature", "bbbb", 2, false),
		hierarchyBranch("fix", "cccc", 3, false),
	}
	revisions := []string{"aaaa", "bbbb", "cccc"}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			buildBranchHierarchyForEachRefArgs(revisions)[1:],
			"refs/heads/main\x00aaaa\x000 0\x000 1\x000 2\n"+
				"refs/heads/feature\x00bbbb\x001 0\x000 0\x000 1\n"+
				"refs/heads/fix\x00cccc\x002 0\x001 0\x000 0\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 41, 0, ""})

	candidates, err := loader.loadDirectBranchCandidates(branches)

	assert.NoError(t, err)
	assert.Equal(t, map[string]int{
		"refs/heads/main":    2,
		"refs/heads/feature": 1,
	}, candidates["refs/heads/fix"])
	runner.CheckForMissingCalls()
}

func BenchmarkProcessLoadedBranchCandidatesFastDenseStack(b *testing.B) {
	const branchCount = 500
	branches := make([]*models.Branch, 0, branchCount)
	relations := make(branchRelationMatrix, branchCount)
	for child := range branchCount {
		name := fmt.Sprintf("branch-%03d", child)
		branches = append(branches, hierarchyBranch(name, fmt.Sprintf("%040x", child+1), int64(child), false))
		candidates := make(map[string]int, child)
		for ancestor := range child {
			candidates[fmt.Sprintf("refs/heads/branch-%03d", ancestor)] = child - ancestor
		}
		relations["refs/heads/"+name] = candidates
	}
	loader := &BranchLoader{GitCommon: &GitCommon{version: &GitVersion{2, 41, 0, ""}}}

	// The dense matrix contains O(n²) strict relations. Modern post-load
	// processing must not add transitive-reduction work on top of that input.
	b.ResetTimer()
	for range b.N {
		candidates, err := loader.processLoadedBranchCandidates(relations, branches)
		if err != nil {
			b.Fatal(err)
		}
		if len(candidates) != branchCount {
			b.Fatalf("got %d child rows, want %d", len(candidates), branchCount)
		}
	}
}

func branchHierarchyLoader(t *testing.T, runner *oscommands.FakeCmdObjRunner, version *GitVersion) *BranchLoader {
	t.Helper()
	gitCommon := buildGitCommon(commonDeps{runner: runner, gitVersion: version})
	return &BranchLoader{Common: gitCommon.Common, GitCommon: gitCommon, cmd: gitCommon.cmd}
}

func TestLoadBranchRelationMatrixLegacy_UsesMergedAndExcludesSameTips(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("grand", "gggg", 1, false),
		hierarchyBranch("parent-a", "aaaa", 2, false),
		hierarchyBranch("parent-b", "bbbb", 3, false),
		hierarchyBranch("child", "cccc", 4, false),
		hierarchyBranch("alias", "cccc", 5, false),
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--merged=gggg", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00gggg\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00gggg\nrefs/heads/parent-a\x00aaaa\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=bbbb", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00gggg\nrefs/heads/parent-b\x00bbbb\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=cccc", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00gggg\nrefs/heads/parent-a\x00aaaa\nrefs/heads/parent-b\x00bbbb\nrefs/heads/child\x00cccc\nrefs/heads/alias\x00cccc\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=cccc", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00gggg\nrefs/heads/parent-a\x00aaaa\nrefs/heads/parent-b\x00bbbb\nrefs/heads/child\x00cccc\nrefs/heads/alias\x00cccc\n", nil)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixLegacy(branches)

	assert.NoError(t, err)
	assert.Equal(t, map[string]int{
		"refs/heads/grand":    unknownBranchDistance,
		"refs/heads/parent-a": unknownBranchDistance,
		"refs/heads/parent-b": unknownBranchDistance,
	}, matrix["refs/heads/child"])
	assert.NotContains(t, matrix["refs/heads/child"], "refs/heads/child")
	assert.NotContains(t, matrix["refs/heads/child"], "refs/heads/alias")
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_LoadsDistancesForMultipleDirectCandidates(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("parent-a", "aaaa", 1, false),
		hierarchyBranch("parent-b", "bbbb", 2, false),
		hierarchyBranch("child", "cccc", 3, false),
	}
	direct := relation("child", map[string]int{
		"refs/heads/parent-a": unknownBranchDistance,
		"refs/heads/parent-b": unknownBranchDistance,
	})
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"rev-list", "--count", "aaaa..cccc"}, "2\n", nil).
		ExpectGitArgs([]string{"rev-list", "--count", "bbbb..cccc"}, "1\n", nil)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	err := loader.loadLegacyBranchDistances(direct, branches)

	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"refs/heads/parent-a": 2, "refs/heads/parent-b": 1}, direct["refs/heads/child"])
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_SkipsDistanceForSoleCandidate(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("parent", "aaaa", 1, false),
		hierarchyBranch("child", "cccc", 2, false),
	}
	direct := relation("child", map[string]int{"refs/heads/parent": unknownBranchDistance})
	runner := oscommands.NewFakeRunner(t)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	err := loader.loadLegacyBranchDistances(direct, branches)

	assert.NoError(t, err)
	assert.Equal(t, unknownBranchDistance, direct["refs/heads/child"]["refs/heads/parent"])
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_RejectsMalformedDistance(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("parent-a", "aaaa", 1, false),
		hierarchyBranch("parent-b", "bbbb", 2, false),
		hierarchyBranch("child", "cccc", 3, false),
	}
	direct := relation("child", map[string]int{
		"refs/heads/parent-a": unknownBranchDistance,
		"refs/heads/parent-b": unknownBranchDistance,
	})
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"rev-list", "--count", "aaaa..cccc"}, "not-a-count\n", nil)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	err := loader.loadLegacyBranchDistances(direct, branches)

	assert.EqualError(t, err, "parse branch distance for refs/heads/parent-a..refs/heads/child: strconv.Atoi: parsing \"not-a-count\": invalid syntax")
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_RejectsMalformedNonRefLines(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaa", 1, false),
		hierarchyBranch("feature", "bbbb", 2, false),
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"},
			"refs/heads/main\x00aaaa\nwarning: malformed output\x00aaaa\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixLegacy(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "malformed branch hierarchy ref: warning: malformed output")
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_IgnoresBlankAndConcurrentRefs(t *testing.T) {
	branches := []*models.Branch{hierarchyBranch("main", "aaaa", 1, false)}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"},
			"\nrefs/heads/concurrent\x00cccc\nrefs/heads/main\x00aaaa\n\n",
			nil,
		).
		ExpectGitArgs([]string{"check-ref-format", "refs/heads/concurrent"}, "", nil)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixLegacy(branches)

	assert.NoError(t, err)
	assert.Empty(t, matrix["refs/heads/main"])
	runner.CheckForMissingCalls()
}

func TestLoadBranchRelationMatrixLegacy_RejectsMalformedConcurrentRefs(t *testing.T) {
	for _, malformedRef := range []string{"refs/heads/bad ref", "refs/heads/.."} {
		t.Run(malformedRef, func(t *testing.T) {
			branches := []*models.Branch{hierarchyBranch("main", "aaaa", 1, false)}
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs(
					[]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"},
					malformedRef+"\x00cccc\nrefs/heads/main\x00aaaa\n",
					nil,
				).
				ExpectGitArgs([]string{"check-ref-format", malformedRef}, "", errors.New("invalid ref"))
			loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

			matrix, err := loader.loadBranchRelationMatrixLegacy(branches)

			assert.Nil(t, matrix)
			assert.EqualError(t, err, "invalid ref")
			runner.CheckForMissingCalls()
		})
	}
}

func TestLoadBranchRelationMatrixLegacy_RejectsSwappedKnownRefs(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("main", "aaaa", 1, false),
		hierarchyBranch("feature", "bbbb", 2, false),
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"},
			"refs/heads/main\x00bbbb\nrefs/heads/feature\x00aaaa\n",
			nil,
		)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	matrix, err := loader.loadBranchRelationMatrixLegacy(branches)

	assert.Nil(t, matrix)
	assert.EqualError(t, err, "branch hierarchy snapshot mismatch for refs/heads/main: expected aaaa, got bbbb")
	runner.CheckForMissingCalls()
}

func TestLoadDirectBranchCandidatesLegacy_UsesIndependentOIDsAndExpandsSameTips(t *testing.T) {
	branches := []*models.Branch{
		hierarchyBranch("grand", "aaaa", 1, false),
		hierarchyBranch("parent-one", "bbbb", 2, false),
		hierarchyBranch("parent-two", "bbbb", 3, false),
		hierarchyBranch("child", "cccc", 4, false),
	}
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--merged=aaaa", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00aaaa\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=bbbb", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00aaaa\nrefs/heads/parent-one\x00bbbb\nrefs/heads/parent-two\x00bbbb\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=bbbb", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00aaaa\nrefs/heads/parent-one\x00bbbb\nrefs/heads/parent-two\x00bbbb\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--merged=cccc", "--format=%(refname)%00%(objectname)", "refs/heads"}, "refs/heads/grand\x00aaaa\nrefs/heads/parent-one\x00bbbb\nrefs/heads/parent-two\x00bbbb\nrefs/heads/child\x00cccc\n", nil).
		ExpectGitArgs([]string{"merge-base", "--independent", "aaaa", "bbbb"}, "bbbb\n", nil).
		ExpectGitArgs([]string{"rev-list", "--count", "bbbb..cccc"}, "1\n", nil).
		ExpectGitArgs([]string{"rev-list", "--count", "bbbb..cccc"}, "1\n", nil)
	loader := branchHierarchyLoader(t, runner, &GitVersion{2, 40, 0, ""})

	candidates, err := loader.loadDirectBranchCandidates(branches)

	assert.NoError(t, err)
	assert.Equal(t, map[string]int{
		"refs/heads/parent-one": 1,
		"refs/heads/parent-two": 1,
	}, candidates["refs/heads/child"])
	runner.CheckForMissingCalls()
}
