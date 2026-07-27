# Local Branch Hierarchy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an optional local-branch hierarchy sort that infers nearest local tip ancestors, orders each subtree by commit date, and renders one compact ASCII dash per downstream level.

**Architecture:** Remove the existing row-zero HEAD assumption first, then extract and generalize the existing ahead/behind query parser. The feature uses a batched ahead/behind matrix on Git 2.41+, a `for-each-ref --merged` fallback on older Git, a pure deterministic forest builder, and a presentation-only depth field.

**Tech Stack:** Go, Git plumbing commands, `stretchr/testify`, lazygit fake command runners, lazygit integration-test drivers, JSON-schema/config documentation generation.

## Global Constraints

- The default local branch sort remains `date`; hierarchy is opt-in as `git.localBranchSortOrder: hierarchy`.
- Parent relationships describe current local branch tips, not historical branch creation or remote tracking.
- Modern ancestry queries use unique commit object IDs, never branch names, inside `%(ahead-behind:...)` atoms.
- Modern `--format` arguments are batched at a maximum of 16 KiB.
- Roots and siblings sort by descending tip commit timestamp and then ascending branch name.
- Parent ties resolve by direct ancestry, distance, unambiguous child count, configured `git.mainBranches`, oldest tip, then branch name.
- Detached and unborn HEAD entries do not participate in ancestry inference and remain standalone first roots.
- Filtered results preserve existing fuzzy relevance ordering and suppress all hierarchy dashes.
- Deep prefixes are shortened before reducing the existing minimum rendered branch-name width.
- Edit only `pkg/i18n/english.go`; never edit `pkg/i18n/translations/`.
- Generate feature documentation under `docs-master/` and `schema-master/`; never edit `docs/` or `schema/`.
- Do not stage or alter the pre-existing `go.sum` worktree change or `.superpowers/brainstorm/` files.
- Do not create a pull request.
- `just` is unavailable in the current shell. Use the exact `justfile` commands listed in this plan. Scope gofumpt to tracked non-vendor Go files so it does not traverse `.worktrees/secondary`.

---

## File Structure

### New files

- `pkg/gui/controllers/helpers/refs_helper_test.go`: pure checked-out-branch finder tests.
- `pkg/commands/git_commands/branch_ahead_behind.go`: shared ahead/behind format construction, parsing, and selection.
- `pkg/commands/git_commands/branch_ahead_behind_test.go`: parser and command-builder tests moved out of the branch loader tests.
- `pkg/commands/git_commands/branch_hierarchy.go`: ancestry loading, batching, parent inference, forest ordering, and fallback.
- `pkg/commands/git_commands/branch_hierarchy_test.go`: pure inference and modern/legacy loader tests.

### Modified source files

- `pkg/gui/controllers/helpers/refs_helper.go`: HEAD lookup, move-commits lookup, sort-menu option.
- `pkg/gui/controllers/helpers/merge_and_rebase_helper.go`: HEAD lookup without creating a constructor cycle.
- `pkg/gui/controllers/helpers/branches_helper.go`: skip the actual HEAD branch during auto-forwarding.
- `pkg/gui/controllers/helpers/refresh_helper.go`: select and reveal the actual visible HEAD row.
- `pkg/gui/controllers/helpers/refresh_helper_test.go`: nonzero and filtered HEAD selection tests.
- `pkg/gui/types/refresh.go`: remove the row-zero selection contract.
- `pkg/commands/models/branch.go`: tip timestamp and hierarchy depth.
- `pkg/commands/git_commands/branch_loader.go`: orchestration, timestamp parsing, hierarchy activation, and flat fallback.
- `pkg/commands/git_commands/branch_loader_test.go`: timestamp parsing and divergence integration tests; ahead/behind helper tests move to the focused test file.
- `pkg/gui/presentation/branches.go`: compact prefix width and neutral rendering.
- `pkg/gui/presentation/branches_test.go`: depth, filtering, color, and truncation cases.
- `pkg/gui/context/branches_context.go`: pass filtering state to presentation.
- `pkg/config/user_config.go`: hierarchy enum and generated-doc comment.
- `pkg/config/user_config_validation.go`: hierarchy validation.
- `pkg/config/user_config_validation_test.go`: valid hierarchy value.
- `pkg/gui/controllers/branches_controller.go`: local sort-menu order.
- `pkg/i18n/english.go`: hierarchy label and ancestry description.
- `pkg/integration/tests/branch/sort_local_branches.go`: hierarchy menu and rendered stack assertions.

### Generated files

- `docs-master/Config.md`: generated local sort-order documentation.
- `schema-master/config.json`: generated hierarchy enum and description.

## Shared Interfaces

Task 1 produces:

```go
func findCheckedOutRef(branches []*models.Branch) (*models.Branch, int, bool)

type branchSelectionContext interface {
	GetItems() []*models.Branch
	SetSelection(int)
}

func selectCheckedOutBranch(context branchSelectionContext) bool
```

Task 2 produces:

```go
type aheadBehind struct {
	baseRevision string
	ahead        int
	behind       int
}

type branchAheadBehind struct {
	refName      string
	aheadBehinds []aheadBehind
}

func parseAheadBehindForEachRefOutput(output string, baseRevisions []string) []branchAheadBehind
func buildAheadBehindForEachRefArgs(baseRevisions []string) []string
func selectBehindForBranch(values []aheadBehind) int
```

Task 3 produces and consumes:

```go
const maxAheadBehindFormatArgBytes = 16 * 1024
const unknownBranchDistance = -1

// Child full ref -> strict ancestor full ref -> commit distance.
type branchRelationMatrix map[string]map[string]int

func batchAheadBehindBaseRevisions(baseRevisions []string, maxFormatArgBytes int) [][]string
func directBranchCandidates(relations branchRelationMatrix) branchRelationMatrix
func inferBranchHierarchy(branches []*models.Branch, candidates branchRelationMatrix, mainBranchNames []string) []*models.Branch

func (self *BranchLoader) loadBranchRelationMatrix(branches []*models.Branch) (branchRelationMatrix, error)
func (self *BranchLoader) loadBranchRelationMatrixFast(branches []*models.Branch) (branchRelationMatrix, error)
func (self *BranchLoader) loadBranchRelationMatrixLegacy(branches []*models.Branch) (branchRelationMatrix, error)
func (self *BranchLoader) loadDirectBranchCandidates(branches []*models.Branch) (branchRelationMatrix, error)
func (self *BranchLoader) loadLegacyBranchDistances(branchesByRef map[string]*models.Branch, candidates branchRelationMatrix) error
func (self *BranchLoader) applyBranchHierarchy(branches []*models.Branch) []*models.Branch
```

---

### Task 1: Remove Checked-Out-Branch Row-Zero Assumptions

**Files:**
- Create: `pkg/gui/controllers/helpers/refs_helper_test.go`
- Modify: `pkg/gui/controllers/helpers/refs_helper.go:193-199,442-445`
- Modify: `pkg/gui/controllers/helpers/merge_and_rebase_helper.go:400-409,519-525`
- Modify: `pkg/gui/controllers/helpers/branches_helper.go:436-458`
- Modify: `pkg/gui/controllers/helpers/refresh_helper.go:1141-1160`
- Modify: `pkg/gui/controllers/helpers/refresh_helper_test.go`
- Modify: `pkg/gui/types/refresh.go:50-64`

**Interfaces:**
- Consumes: `models.Branch.Head`, which is the authoritative checked-out marker.
- Produces: `findCheckedOutRef`, `selectCheckedOutBranch`, and a `RefsHelper.GetCheckedOutRef` implementation independent of list order.

- [ ] **Step 1: Add failing checked-out-ref finder tests**

Create `refs_helper_test.go` with a table that proves position is irrelevant:

```go
func TestFindCheckedOutRef(t *testing.T) {
	testCases := []struct {
		name          string
		branches      []*models.Branch
		expectedBranch *models.Branch
		expectedIndex int
		expectedFound bool
	}{
		{name: "empty", expectedIndex: -1},
		{
			name: "head at zero",
			branches: []*models.Branch{{Name: "main", Head: true}, {Name: "feature"}},
			expectedBranch: &models.Branch{Name: "main", Head: true},
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name: "nested head",
			branches: []*models.Branch{{Name: "main"}, {Name: "feature", Head: true}},
			expectedBranch: &models.Branch{Name: "feature", Head: true},
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name: "no head marker",
			branches: []*models.Branch{{Name: "main"}},
			expectedIndex: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			branch, index, found := findCheckedOutRef(testCase.branches)
			assert.Equal(t, testCase.expectedBranch, branch)
			assert.Equal(t, testCase.expectedIndex, index)
			assert.Equal(t, testCase.expectedFound, found)
		})
	}
}
```

- [ ] **Step 2: Run the finder test and verify RED**

Run:

```bash
go test ./pkg/gui/controllers/helpers -run TestFindCheckedOutRef -count=1
go test ./pkg/gui/controllers/helpers -short
```

Expected: build failure because `findCheckedOutRef` is undefined.

- [ ] **Step 3: Implement the shared finder and replace direct lookups**

Add the package-level helper in `refs_helper.go` so `MergeAndRebaseHelper` can use it without depending on `RefsHelper`, which already depends on the merge helper:

```go
func findCheckedOutRef(branches []*models.Branch) (*models.Branch, int, bool) {
	return lo.FindIndexOf(branches, func(branch *models.Branch) bool {
		return branch.Head
	})
}

func (self *RefsHelper) GetCheckedOutRef() *models.Branch {
	branch, _, _ := findCheckedOutRef(self.c.Model().Branches)
	return branch
}
```

Use `self.GetCheckedOutRef()` in `MoveCommitsToNewBranch`. Use `findCheckedOutRef(self.c.Model().Branches)` in `RebaseOntoRef` and `MergeRefIntoCheckedOutBranch`; preserve current detached/no-branch error handling rather than dereferencing nil. In `AutoForwardBranches`, iterate all branches and skip `branch.Head` instead of slicing `branches[1:]`:

```go
for _, branch := range branches {
	if branch.Head {
		continue
	}
	// Existing remote/worktree/config checks remain unchanged.
}
```

- [ ] **Step 4: Run the finder test and helper package tests**

Run:

```bash
go test ./pkg/gui/controllers/helpers -run TestFindCheckedOutRef -count=1
```

Expected: PASS.

- [ ] **Step 5: Add failing visible-HEAD selection tests**

In `refresh_helper_test.go`, use the real filtered list model so visible indexes and range cancellation are exercised:

```go
func TestSelectCheckedOutBranch(t *testing.T) {
	branches := []*models.Branch{
		{Name: "main"},
		{Name: "feature", Head: true},
	}
	viewModel := context.NewFilteredListViewModel(
		func() []*models.Branch { return branches },
		func(branch *models.Branch) []string { return []string{branch.Name} },
	)
	viewModel.SetSelectionRangeAndMode(0, 1, traits.RangeSelectModeSticky)

	found := selectCheckedOutBranch(viewModel)

	assert.True(t, found)
	assert.Equal(t, 1, viewModel.GetSelectedLineIdx())
	_, _, mode := viewModel.GetSelectionRangeAndMode()
	assert.Equal(t, traits.RangeSelectModeNone, mode)
}

func TestSelectCheckedOutBranchWhenFilteredOut(t *testing.T) {
	branches := []*models.Branch{{Name: "main"}, {Name: "feature", Head: true}}
	viewModel := context.NewFilteredListViewModel(
		func() []*models.Branch { return branches },
		func(branch *models.Branch) []string { return []string{branch.Name} },
	)
	viewModel.SetFilter("main", false)

	assert.False(t, selectCheckedOutBranch(viewModel))
	assert.Equal(t, "main", viewModel.GetSelected().Name)
}
```

Add the required `context` import; do not create a fake list cursor.

- [ ] **Step 6: Run selection tests and verify RED**

Run:

```bash
go test ./pkg/gui/controllers/helpers -run TestSelectCheckedOutBranch -count=1
```

Expected: build failure because `selectCheckedOutBranch` is undefined.

- [ ] **Step 7: Implement nonzero selection and post-render focus**

Add the narrow interface and pure selector in `refresh_helper.go`:

```go
type branchSelectionContext interface {
	GetItems() []*models.Branch
	SetSelection(int)
}

func selectCheckedOutBranch(context branchSelectionContext) bool {
	_, index, found := findCheckedOutRef(context.GetItems())
	if found {
		context.SetSelection(index)
	}
	return found
}
```

In the `SelectCheckedOutBranch` switch case:

1. Reapply the branch filter before reading `GetItems()`.
2. Call `selectCheckedOutBranch`.
3. If it returns true, enqueue a generation-guarded nested UI bounce that calls `self.c.Contexts().Branches.FocusLine(true)` after rendering, following the existing HEAD-commit selection pattern in the same file.
4. Remove `SetOriginY(0)`.
5. If the active filter excludes HEAD, leave selection and filter unchanged.

Update `types/refresh.go` to say "Select the branch whose `Head` marker is true" rather than "the one at the top."

- [ ] **Step 8: Run focused and full unit tests**

Run:

```bash
go test ./pkg/gui/controllers/helpers -run 'Test(FindCheckedOutRef|SelectCheckedOutBranch)' -count=1
go test ./... -short
```

Expected: PASS.

- [ ] **Step 9: Format, lint, build, inspect, and commit Task 1**

Run:

```bash
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l -w
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l
./scripts/golangci-lint-shim.sh run
go build -gcflags='all=-N -l'
git status --short
git diff
git log --oneline -10
```

Expected: no primary tracked Go file from the second gofumpt command, `0 issues.`, successful build, and only Task 1 files plus unrelated pre-existing changes in status.

Commit only Task 1 files:

```bash
git add pkg/gui/controllers/helpers/refs_helper.go pkg/gui/controllers/helpers/refs_helper_test.go pkg/gui/controllers/helpers/merge_and_rebase_helper.go pkg/gui/controllers/helpers/branches_helper.go pkg/gui/controllers/helpers/refresh_helper.go pkg/gui/controllers/helpers/refresh_helper_test.go pkg/gui/types/refresh.go
git commit -m "Find the checked-out branch independently of display order"
```

---

### Task 2: Generalize And Extract Ahead/Behind Parsing

**Files:**
- Create: `pkg/commands/git_commands/branch_ahead_behind.go`
- Create: `pkg/commands/git_commands/branch_ahead_behind_test.go`
- Modify: `pkg/commands/git_commands/branch_loader.go:209-329`
- Modify: `pkg/commands/git_commands/branch_loader_test.go:129-408`

**Interfaces:**
- Consumes: existing Git 2.41 `%(ahead-behind:<revision>)` output.
- Produces: base-revision-aware parsing and query construction used by both divergence and Task 3 hierarchy loading.

- [ ] **Step 1: Move existing helper tests into the focused test file**

Move `TestParseAheadBehindForEachRefOutput`, `TestSelectBehindForBranch`, and `TestBuildAheadBehindForEachRefArgs` from `branch_loader_test.go` into `branch_ahead_behind_test.go`. Initially move the current helper types and functions unchanged into `branch_ahead_behind.go`, and confirm package tests still pass before changing signatures.

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(ParseAheadBehindForEachRefOutput|SelectBehindForBranch|BuildAheadBehindForEachRefArgs)' -count=1
```

Expected: PASS, proving the extraction alone is behavior-preserving.

- [ ] **Step 2: Change parser tests to require base identity and safe invalid fields**

Change test input from `numBases int` to `baseRevisions []string`. Add this exact malformed-middle case:

```go
{
	testName:     "malformed middle field retains later base identity",
	input:        "refs/heads/child\x002 0\x00bad\x007 1\n",
	baseRevisions: []string{"oid-a", "oid-b", "oid-c"},
	expected: []branchAheadBehind{
		{
			refName: "refs/heads/child",
			aheadBehinds: []aheadBehind{
				{baseRevision: "oid-a", ahead: 2, behind: 0},
				{baseRevision: "oid-c", ahead: 7, behind: 1},
			},
		},
	},
}
```

Retain cases for blank input, wrong row width, unreachable/empty bases, malformed numeric fields, and valid multiple rows. Add `baseRevision` to every expected valid field. Keep the empty `selectBehindForBranch` cases expecting zero.

- [ ] **Step 3: Run parser tests and verify RED**

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(ParseAheadBehindForEachRefOutput|SelectBehindForBranch)' -count=1
```

Expected: build failures because the parser still accepts an integer and `aheadBehind` has no `baseRevision`.

- [ ] **Step 4: Implement base-aware parsing and empty selection**

Implement the Task 2 shared interfaces exactly:

```go
type aheadBehind struct {
	baseRevision string
	ahead        int
	behind       int
}

func parseAheadBehindForEachRefOutput(output string, baseRevisions []string) []branchAheadBehind {
	if output == "" {
		return nil
	}

	result := make([]branchAheadBehind, 0)
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		columns := strings.Split(line, "\x00")
		if len(columns) != len(baseRevisions)+1 {
			continue
		}

		values := make([]aheadBehind, 0, len(baseRevisions))
		for index, column := range columns[1:] {
			value, valid := parseAheadBehindField(column)
			if valid {
				value.baseRevision = baseRevisions[index]
				values = append(values, value)
			}
		}
		result = append(result, branchAheadBehind{refName: columns[0], aheadBehinds: values})
	}
	return result
}

func selectBehindForBranch(values []aheadBehind) int {
	if len(values) == 0 {
		return 0
	}
	return lo.MinBy(values, func(a, b aheadBehind) bool { return a.ahead < b.ahead }).behind
}
```

Keep valid values in source-column order so equal-distance divergence ties retain configured main-branch order. Change `getBehindBaseBranchValuesFast` to pass `mainBranchRefs` instead of their count.

- [ ] **Step 5: Add a divergence test for an all-invalid row**

Extend the fast-path test table or add a focused test whose fake Git output is:

```text
refs/heads/feat-x\x00not-a-count
```

Initialize `BehindBaseBranch` to a nonzero stale value and assert the fast loader stores zero without shifting or failing:

```go
assert.NoError(t, err)
assert.Equal(t, int32(0), feat.BehindBaseBranch.Load())
```

- [ ] **Step 6: Run focused and full package tests**

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(ParseAheadBehindForEachRefOutput|SelectBehindForBranch|BuildAheadBehindForEachRefArgs|GetBehindBaseBranchValuesForAllBranches)' -count=1
go test ./pkg/commands/git_commands -short
```

Expected: PASS.

- [ ] **Step 7: Format, verify, inspect, and commit Task 2**

Run:

```bash
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l -w
go test ./... -short
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l
./scripts/golangci-lint-shim.sh run
go build -gcflags='all=-N -l'
git status --short
git diff
git log --oneline -10
```

Expected: all checks pass and only Task 2 files plus unrelated pre-existing changes appear.

Commit only Task 2 files:

```bash
git add pkg/commands/git_commands/branch_ahead_behind.go pkg/commands/git_commands/branch_ahead_behind_test.go pkg/commands/git_commands/branch_loader.go pkg/commands/git_commands/branch_loader_test.go
git commit -m "Share ahead-behind results without losing base identity"
```

---

### Task 3: Add Local Branch Hierarchy Behavior

**Files:**
- Create: `pkg/commands/git_commands/branch_hierarchy.go`
- Create: `pkg/commands/git_commands/branch_hierarchy_test.go`
- Modify: `pkg/commands/models/branch.go`
- Modify: `pkg/commands/git_commands/branch_loader.go`
- Modify: `pkg/commands/git_commands/branch_loader_test.go`
- Modify: `pkg/gui/presentation/branches.go`
- Modify: `pkg/gui/presentation/branches_test.go`
- Modify: `pkg/gui/context/branches_context.go`
- Modify: `pkg/config/user_config.go`
- Modify: `pkg/config/user_config_validation.go`
- Modify: `pkg/config/user_config_validation_test.go`
- Modify: `pkg/gui/controllers/helpers/refs_helper.go`
- Modify: `pkg/gui/controllers/branches_controller.go`
- Modify: `pkg/i18n/english.go`
- Modify: `pkg/integration/tests/branch/sort_local_branches.go`
- Generate: `docs-master/Config.md`
- Generate: `schema-master/config.json`

**Interfaces:**
- Consumes: Task 1 HEAD lookup behavior and Task 2 ahead/behind parser/builder.
- Produces: hierarchy sorting, `models.Branch.HierarchyDepth`, `hierarchy` config/menu option, compact display, and end-to-end coverage.

- [ ] **Step 1: Add failing model and timestamp parsing expectations**

Add to `models.Branch` test expectations before changing the model. In every valid `TestObtainBranch` timestamp case, require:

```go
CommitUnixTimestamp: unixTimestamp,
```

where the fixture computes both forms once:

```go
unixTimestamp := now - int64(2.5*60*60)
timeStamp := strconv.FormatInt(unixTimestamp, 10)
```

Add a malformed timestamp scenario expecting `CommitUnixTimestamp == 0` and empty `Recency`.

Run:

```bash
go test ./pkg/commands/git_commands -run TestObtainBranch -count=1
```

Expected: build failure because `models.Branch.CommitUnixTimestamp` does not exist.

- [ ] **Step 2: Add model fields and parse the timestamp once**

Add documented fields in `pkg/commands/models/branch.go`:

```go
// CommitUnixTimestamp is the committer date of the branch tip.
CommitUnixTimestamp int64
// HierarchyDepth is the inferred local-branch ancestry depth in hierarchy mode.
HierarchyDepth int
```

In `obtainBranch`, parse `split[7]` once, store it when valid, and derive `Recency` from that value only when `storeCommitDateAsRecency` is true. Run `TestObtainBranch` and expect PASS.

- [ ] **Step 3: Add failing pure hierarchy inference tests**

Create `branch_hierarchy_test.go` with a branch helper and full-ref relation helper:

```go
func hierarchyBranch(name, hash string, timestamp int64, head bool) *models.Branch {
	return &models.Branch{
		Name: name, CommitHash: hash, CommitUnixTimestamp: timestamp, Head: head,
	}
}

func relation(child string, candidates map[string]int) branchRelationMatrix {
	return branchRelationMatrix{"refs/heads/" + child: candidates}
}
```

Write separate table/subtests with explicit expected `name:depth` output:

```go
func hierarchyShape(branches []*models.Branch) []string {
	return lo.Map(branches, func(branch *models.Branch, _ int) string {
		return fmt.Sprintf("%s:%d", branch.Name, branch.HierarchyDepth)
	})
}
```

Required cases:

```text
linear:              main:0, feature:1, fix:2
siblings by date:    main:0, newer:1, older:1
unrelated roots:     newer-root:0, older-root:0
same-tip refs:       main:0, alias:0
nested HEAD subtree: old-root:0, child-head:1, newer-unrelated-root:0
detached HEAD:       detached:0, main:0, feature:1
unborn HEAD:         unborn:0
```

For transitive reduction, provide `main -> feature`, `main -> fix`, and `feature -> fix`; assert `fix` has depth 2. Add equal-distance tie tests proving, in order:

1. Candidate with more unambiguous children wins.
2. Configured `mainBranchNames` candidate wins.
3. Older candidate timestamp wins.
4. Alphabetically smaller candidate wins.

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(DirectBranchCandidates|InferBranchHierarchy)' -count=1
```

Expected: build failure because relation types and inference functions are undefined.

- [ ] **Step 4: Implement pure candidate reduction and forest ordering**

In `branch_hierarchy.go`:

1. Key all real branches by `FullRefName()`.
2. `directBranchCandidates` removes candidate A for child C when another candidate B for C has A in `relations[B]`.
3. Treat the supplied candidates as already direct. Keep a sole candidate regardless of distance; when several remain, keep only the minimum nonnegative distance.
4. Record unique survivors and count their unambiguous children before resolving remaining ties.
5. Resolve ties with the exact Global Constraints order.
6. Build local `childrenByParent` and `parentByChild` maps.
7. Sort roots and each child slice by timestamp descending, then name ascending.
8. Move the root containing attached HEAD to the first root position.
9. Prepend synthetic `Head` entries with empty `CommitHash` as depth-zero roots.
10. Flatten depth-first and assign `HierarchyDepth` during traversal.

Use local maps only; do not add parent or child pointers to `models.Branch`. Add an input-permutation test to prove tie resolution is deterministic. Run the focused inference tests and expect PASS.

- [ ] **Step 5: Add failing batch and modern matrix tests**

Test `batchAheadBehindBaseRevisions` with SHA-1 and SHA-256-length strings. For every returned batch, reconstruct the format argument using `buildAheadBehindForEachRefArgs(batch)[2]` and assert:

```go
assert.LessOrEqual(t, len(formatArg), maxAheadBehindFormatArgBytes)
```

Assert no revision is lost, duplicated, or split. Add a modern loader fake-runner test with:

```go
branches := []*models.Branch{
	hierarchyBranch("main)", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, false),
	hierarchyBranch("feature", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 2, true),
}
```

Expect a format atom containing only the safe `a...` and `b...` OIDs. Return rows where feature is ahead 1/behind 0 of main and assert:

```go
assert.Equal(t, 1, matrix["refs/heads/feature"]["refs/heads/main)"])
```

Also test that duplicate branch tips create one format atom but expand into two candidate refs for a later child.

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(BatchAheadBehindBaseRevisions|LoadBranchRelationMatrixFast)' -count=1
```

Expected: build failure because batching and fast loading are undefined.

- [ ] **Step 6: Implement safe modern batching and relation expansion**

Build each format-size calculation from:

```go
const formatPrefix = "--format=%(refname)"
const atomPrefix = "%00%(ahead-behind:"
const atomSuffix = ")"
```

Deduplicate nonempty branch `CommitHash` values while retaining `map[string][]*models.Branch`. For each batch:

1. Run the Task 2 `buildAheadBehindForEachRefArgs` command.
2. Parse with the same batch revisions.
3. Resolve each row's child full ref.
4. Keep fields with `behind == 0` and `baseRevision != child.CommitHash`.
5. Expand the OID to all branch refs at that tip and store `ahead`.
6. On any command error, discard the partial matrix and return the error.

Set the version dispatch in `loadBranchRelationMatrix` to `self.version.IsAtLeast(2, 41, 0)`. Run the modern tests and expect PASS.

- [ ] **Step 7: Add failing legacy relation and distance tests**

Use `GitVersion{2, 40, 0, ""}` and fake expectations for:

```text
git for-each-ref --merged=<childOID> --format=%(refname) refs/heads
```

Return full refs including child, same-tip alias, a direct parent, and a hidden grandparent. Assert child and same-tip refs are excluded. After direct reduction leaves two merge-side candidates, expect only:

```text
git rev-list --count <candidateOID>..<childOID>
```

for those two candidates. Add another test with one direct candidate and assert no rev-list expectation is needed.

Run:

```bash
go test ./pkg/commands/git_commands -run TestLoadBranchRelationMatrixLegacy -count=1
```

Expected: failure because the legacy implementation is absent.

- [ ] **Step 8: Implement the older-Git path**

For each real local child, run the exact `for-each-ref --merged=<childOID>` command and parse newline-delimited full refs. `loadBranchRelationMatrixLegacy` returns the complete strict relation with `unknownBranchDistance`. `loadDirectBranchCandidates` applies `directBranchCandidates` to either version's matrix; on old Git it then calls `loadLegacyBranchDistances` only for children with multiple direct candidates. Parse trimmed counts with `strconv.Atoi`; return an error for command failures or malformed counts. Keep the fake runner checks deterministic and call `runner.CheckForMissingCalls()` in every loader test.

Run modern, legacy, and pure inference tests together and expect PASS.

- [ ] **Step 9: Add failing BranchLoader hierarchy and fallback tests**

Add tests that exercise hierarchy orchestration rather than only pure helpers:

1. Git 2.41 hierarchy mode returns a nested attached HEAD instead of moving it to row zero.
2. Date, recency, and alphabetical modes still move HEAD to row zero.
3. Detached and unborn synthetic HEAD entries never appear in matrix base atoms.
4. A matrix command error returns timestamp/name-sorted flat branches, resets all depths to zero, and promotes HEAD to row zero.

Expected hierarchy success shape:

```text
main:0
feature:1
head-fix:2
```

Expected fallback shape:

```text
head-fix:0
feature:0
main:0
```

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(BranchLoaderHierarchy|ApplyBranchHierarchy)' -count=1
```

Expected: failure because `BranchLoader.Load` does not activate hierarchy.

- [ ] **Step 10: Integrate hierarchy mode and flat fallback into BranchLoader**

In `Load`:

1. Keep initial `getRawBranches` sorting for `hierarchy` as `-committerdate`.
2. Set real/synthetic HEAD recency exactly as today.
3. Move real HEAD to index zero only when sort order is not `hierarchy`.
4. Populate config/upstream and stale divergence fields as today.
5. If sort order is `hierarchy`, call `applyBranchHierarchy` before scheduling base-divergence loading.

`applyBranchHierarchy` must zero all depths, call `loadDirectBranchCandidates`, and pass its direct, distance-populated candidates to `inferBranchHierarchy`. On error, log once, sort timestamp descending/name ascending, move HEAD to zero, and return without propagating the error so the panel remains usable.

Run all `pkg/commands/git_commands` tests and expect PASS.

- [ ] **Step 11: Add failing config, menu, and English-string tests/expectations**

In `TestUserConfigValidate_enums`, add:

```go
{value: "hierarchy", valid: true},
```

Update `sort_local_branches.go` menu expectations to include:

```go
Contains("h ( ) Hierarchy"),
```

Add English fields and values alongside existing sort strings:

```go
SortByHierarchy           string
SortBasedOnBranchAncestry string

SortByHierarchy:           "Hierarchy",
SortBasedOnBranchAncestry: "(based on branch ancestry)",
```

Run config tests before source enum changes:

```bash
go test ./pkg/config -run TestUserConfigValidate_enums -count=1
```

Expected: FAIL because validation rejects `hierarchy`.

- [ ] **Step 12: Implement config and sort-menu exposure**

Update `LocalBranchSortOrder` without hard-wrapping a sentence in its generated doc comment:

```go
// How branches are sorted in the local branches view.
// One of: 'date' (default) | 'recency' | 'alphabetical' | 'hierarchy'
// Can be changed from within Lazygit with the Sort Order menu (`s`) in the branches panel.
LocalBranchSortOrder string `yaml:"localBranchSortOrder" jsonschema:"enum=date,enum=recency,enum=alphabetical,enum=hierarchy"`
```

Add `hierarchy` to validation. In `CreateSortOrderMenu`, add:

```go
"hierarchy": {
	label: self.c.Tr.SortByHierarchy,
	description: self.c.Tr.SortBasedOnBranchAncestry,
	keys: menuKey('h'),
},
```

Pass `[]string{"recency", "alphabetical", "date", "hierarchy"}` from the local branches controller. Do not add hierarchy to the remote controller's explicit options.

Run config tests and `go test ./pkg/gui/controllers/... -short`; expect PASS.

- [ ] **Step 13: Add failing hierarchy presentation tests**

Add `isFiltering bool` to the presentation test table before changing production signatures. Required expected name columns include:

```text
depth 0, wide:               branch_name
depth 1, wide:               -branch_name
depth 3, wide:               ---branch_name
depth 3, filtering:          branch_name
depth 2, width 14:           --branch_…
depth 10, very narrow long:  bra…
depth 10, narrow abc:        -abc
```

Add a list-level test with input order `beta`, `alpha`, both nonzero depth, and `isFiltering=true`; assert output stays `beta`, `alpha` with no prefixes. With color enabled, assert the prefix is outside the branch style:

```go
assert.Equal(t,
	style.FgDefault.Sprint("--")+GetBranchTextStyle("branch_name").Sprint("branch_name"),
	actualNameColumn,
)
```

Run:

```bash
go test ./pkg/gui/presentation -run 'Test_(getBranchDisplayStrings|GetBranchListDisplayStrings)' -count=1
```

Expected: build failures because presentation has no filtering argument or depth handling.

- [ ] **Step 14: Implement compact depth rendering and filter suppression**

Thread `isFiltering bool` through `GetBranchListDisplayStrings` and `getBranchDisplayStrings`. Pass `viewModel.IsFiltering()` from `branches_context.go`.

After status/worktree width deductions and before name truncation, calculate:

```go
nameWidth := utils.StringWidth(displayName)
minimumNameWidth := nameWidth
if nameWidth > 3 {
	minimumNameWidth = 4 // three characters plus the ellipsis
}
depth := lo.Ternary(isFiltering, 0, b.HierarchyDepth)
visibleDepth := min(depth, max(availableWidth-minimumNameWidth, 0))
prefix := strings.Repeat("-", visibleDepth)
availableWidth -= visibleDepth
```

Keep the existing branch-name truncation after this deduction. Construct the final name as:

```go
coloredName := style.FgDefault.Sprint(prefix) + nameTextStyle.Sprint(displayName)
```

Then append worktree/status/divergence exactly as today so their width and padding calculations include the visible prefix. Run all presentation tests and expect PASS.

- [ ] **Step 15: Complete the hierarchy integration scenario**

After the existing alphabetical assertion in `SortLocalBranches`, reopen the sort menu, assert Alphabetical is selected and Hierarchy is available, choose Hierarchy, and assert:

```go
t.Views().Branches().
	IsFocused().
	Lines(
		Contains("master").IsSelected(),
		Contains("-first"),
		Contains("--second"),
		Contains("---third"),
	)
```

Update the test description to "Sort local branches by recency, date, alphabetically, or hierarchy." Add the hierarchy row to every sort-menu assertion in this test.

- [ ] **Step 16: Generate config documentation and inspect generated scope**

Run the unavailable `just generate` recipe directly:

```bash
go generate ./...
```

Inspect:

```bash
git status --short
git diff -- docs-master/Config.md schema-master/config.json pkg/integration/tests/test_list.go
```

Expected: `docs-master/Config.md` and `schema-master/config.json` describe the new enum while retaining `date` as default. `pkg/integration/tests/test_list.go` remains unchanged because no test was added or renamed. Do not stage any generated file outside the expected scope without investigating it first.

- [ ] **Step 17: Run focused RED-GREEN verification for Task 3**

Run:

```bash
go test ./pkg/commands/git_commands -run 'Test(ObtainBranch|DirectBranchCandidates|InferBranchHierarchy|BatchAheadBehindBaseRevisions|LoadBranchRelationMatrix|BranchLoaderHierarchy|ApplyBranchHierarchy)' -count=1
go test ./pkg/gui/presentation -run 'Test_(getBranchDisplayStrings|GetBranchListDisplayStrings)' -count=1
go test ./pkg/config -run TestUserConfigValidate_enums -count=1
go test ./pkg/gui/controllers/... -short
```

Expected: PASS.

- [ ] **Step 18: Run the full required verification**

Run in this order:

```bash
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l -w
go test ./... -short
git ls-files -z '*.go' ':!vendor/**' | xargs -0 go tool gofumpt -l
./scripts/golangci-lint-shim.sh run
go build -gcflags='all=-N -l'
go test -timeout 30m pkg/integration/clients/*.go -run 'TestIntegration/branch/sort_local_branches'
```

Expected: unit tests PASS, no primary tracked Go file from the check-only gofumpt command, `0 issues.`, successful build, and the targeted integration test PASS.

- [ ] **Step 19: Review and commit Task 3 without unrelated files**

Inspect before staging:

```bash
git status --short
git diff
git log --oneline -10
```

Stage only the feature files:

```bash
git add pkg/commands/models/branch.go pkg/commands/git_commands/branch_loader.go pkg/commands/git_commands/branch_loader_test.go pkg/commands/git_commands/branch_hierarchy.go pkg/commands/git_commands/branch_hierarchy_test.go pkg/gui/presentation/branches.go pkg/gui/presentation/branches_test.go pkg/gui/context/branches_context.go pkg/config/user_config.go pkg/config/user_config_validation.go pkg/config/user_config_validation_test.go pkg/gui/controllers/helpers/refs_helper.go pkg/gui/controllers/branches_controller.go pkg/i18n/english.go pkg/integration/tests/branch/sort_local_branches.go docs-master/Config.md schema-master/config.json
git diff --cached --check
git diff --cached --stat
```

Confirm `go.sum` and `.superpowers/brainstorm/` are not staged. Commit with:

```bash
git commit -m "Make stacked local branches readable at a glance"
```

---

## Final Review

- [ ] Confirm the three implementation commits follow the design commit `655eca7a4` and each is independently buildable.
- [ ] Confirm `git status --short` contains only the pre-existing `go.sum` change and `.superpowers/brainstorm/` files.
- [ ] Confirm `git log --oneline -4` shows the design commit followed by the three implementation commits.
- [ ] Re-run the full verification command block from Task 3 Step 18 if any file changed after it.
- [ ] Request a code review using the `requesting-code-review` skill; address findings with separate `fixup!` commits targeting the relevant implementation commit.
- [ ] Do not squash fixup commits, amend existing commits, push, or create a pull request.
