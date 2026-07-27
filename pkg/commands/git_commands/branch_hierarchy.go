package git_commands

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type branchRelationMatrix map[string]map[string]int

const unknownBranchDistance = -1

const (
	maxAheadBehindFormatArgBytes = 16 * 1024
	formatPrefix                 = "--format=%(refname)%00%(objectname)"
	atomPrefix                   = "%00%(ahead-behind:"
	atomSuffix                   = ")"
)

func batchAheadBehindBaseRevisions(baseRevisions []string) [][]string {
	if len(baseRevisions) == 0 {
		return nil
	}

	batches := make([][]string, 0, 1)
	batch := make([]string, 0)
	formatLength := len(formatPrefix)
	for _, revision := range baseRevisions {
		atomLength := len(atomPrefix) + len(revision) + len(atomSuffix)
		if len(batch) > 0 && formatLength+atomLength > maxAheadBehindFormatArgBytes {
			batches = append(batches, batch)
			batch = make([]string, 0)
			formatLength = len(formatPrefix)
		}
		batch = append(batch, revision)
		formatLength += atomLength
	}
	return append(batches, batch)
}

func (self *BranchLoader) loadBranchRelationMatrix(branches []*models.Branch) (branchRelationMatrix, error) {
	if self.version.IsAtLeast(2, 41, 0) {
		return self.loadBranchRelationMatrixFast(branches)
	}
	return self.loadBranchRelationMatrixLegacy(branches)
}

func (self *BranchLoader) loadBranchRelationMatrixFast(branches []*models.Branch) (branchRelationMatrix, error) {
	branchesByCommitHash := map[string][]*models.Branch{}
	baseRevisions := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch.CommitHash == "" {
			continue
		}
		if _, exists := branchesByCommitHash[branch.CommitHash]; !exists {
			baseRevisions = append(baseRevisions, branch.CommitHash)
		}
		branchesByCommitHash[branch.CommitHash] = append(branchesByCommitHash[branch.CommitHash], branch)
	}

	branchByRef := make(map[string]*models.Branch, len(branches))
	expectedRefs := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch.CommitHash != "" {
			ref := branch.FullRefName()
			branchByRef[ref] = branch
			expectedRefs = append(expectedRefs, ref)
		}
	}

	matrix := branchRelationMatrix{}
	for _, batch := range batchAheadBehindBaseRevisions(baseRevisions) {
		output, err := self.cmd.New(buildBranchHierarchyForEachRefArgs(batch)).DontLog().RunWithOutput()
		if err != nil {
			return nil, err
		}
		rows, err := parseCompleteAheadBehindOutput(output, batch, branchByRef, expectedRefs)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			child := branchByRef[row.refName]
			if child == nil {
				continue
			}
			if matrix[row.refName] == nil {
				matrix[row.refName] = map[string]int{}
			}
			for _, value := range row.aheadBehinds {
				if value.behind != 0 || value.baseRevision == child.CommitHash {
					continue
				}
				for _, candidate := range branchesByCommitHash[value.baseRevision] {
					matrix[row.refName][candidate.FullRefName()] = value.ahead
				}
			}
		}
	}
	return matrix, nil
}

func buildBranchHierarchyForEachRefArgs(baseRevisions []string) []string {
	format := formatPrefix
	for _, revision := range baseRevisions {
		format += atomPrefix + revision + atomSuffix
	}
	return NewGitCmd("for-each-ref").
		Arg(format).
		Arg("refs/heads").
		ToArgv()
}

func parseCompleteAheadBehindOutput(
	output string,
	baseRevisions []string,
	branchByRef map[string]*models.Branch,
	expectedRefs []string,
) ([]branchAheadBehind, error) {
	expectedRefSet := make(map[string]struct{}, len(expectedRefs))
	for _, ref := range expectedRefs {
		expectedRefSet[ref] = struct{}{}
	}
	seenRefs := make(map[string]struct{}, len(expectedRefs))
	rows := make([]branchAheadBehind, 0, len(expectedRefs))
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		columns := strings.Split(line, "\x00")
		if len(columns) != len(baseRevisions)+2 {
			return nil, fmt.Errorf(
				"malformed branch hierarchy row for %s: expected %d columns, got %d",
				columns[0], len(baseRevisions)+2, len(columns),
			)
		}
		branch := branchByRef[columns[0]]
		if branch == nil {
			continue
		}
		if columns[1] != branch.CommitHash {
			return nil, fmt.Errorf(
				"branch hierarchy snapshot mismatch for %s: expected %s, got %s",
				columns[0], branch.CommitHash, columns[1],
			)
		}
		if _, expected := expectedRefSet[columns[0]]; expected {
			seenRefs[columns[0]] = struct{}{}
		}
		values := make([]aheadBehind, 0, len(baseRevisions))
		for index, column := range columns[2:] {
			value, valid := parseAheadBehindField(column)
			if valid {
				value.baseRevision = baseRevisions[index]
				values = append(values, value)
			}
		}
		rows = append(rows, branchAheadBehind{refName: columns[0], aheadBehinds: values})
	}
	for _, ref := range expectedRefs {
		if _, seen := seenRefs[ref]; !seen {
			return nil, fmt.Errorf("branch hierarchy output missing row for %s", ref)
		}
	}
	return rows, nil
}

func (self *BranchLoader) loadBranchRelationMatrixLegacy(branches []*models.Branch) (branchRelationMatrix, error) {
	branchByRef := make(map[string]*models.Branch, len(branches))
	for _, branch := range branches {
		if branch.CommitHash != "" {
			branchByRef[branch.FullRefName()] = branch
		}
	}

	matrix := make(branchRelationMatrix, len(branchByRef))
	validatedConcurrentRefs := map[string]struct{}{}
	for _, child := range branches {
		if child.CommitHash == "" {
			continue
		}
		output, err := self.cmd.New(
			NewGitCmd("for-each-ref").
				Arg("--merged=" + child.CommitHash).
				Arg("--format=%(refname)%00%(objectname)").
				Arg("refs/heads").
				ToArgv(),
		).DontLog().RunWithOutput()
		if err != nil {
			return nil, err
		}

		childRef := child.FullRefName()
		candidates := map[string]int{}
		childSeen := false
		for _, line := range strings.Split(output, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			columns := strings.Split(line, "\x00")
			if len(columns) != 2 {
				return nil, fmt.Errorf("malformed branch hierarchy row: %s", line)
			}
			candidateRef := strings.TrimSpace(columns[0])
			candidateOID := strings.TrimSpace(columns[1])
			if !strings.HasPrefix(candidateRef, "refs/heads/") || candidateRef == "refs/heads/" {
				return nil, fmt.Errorf("malformed branch hierarchy ref: %s", candidateRef)
			}
			candidate := branchByRef[candidateRef]
			if candidate == nil {
				if _, validated := validatedConcurrentRefs[candidateRef]; !validated {
					err := self.cmd.New(
						NewGitCmd("check-ref-format").Arg(candidateRef).ToArgv(),
					).DontLog().Run()
					if err != nil {
						return nil, err
					}
					validatedConcurrentRefs[candidateRef] = struct{}{}
				}
				continue
			}
			if candidateOID != candidate.CommitHash {
				return nil, fmt.Errorf(
					"branch hierarchy snapshot mismatch for %s: expected %s, got %s",
					candidateRef, candidate.CommitHash, candidateOID,
				)
			}
			if candidateRef == childRef {
				childSeen = true
			}
			if candidate.CommitHash == child.CommitHash {
				continue
			}
			candidates[candidateRef] = unknownBranchDistance
		}
		if !childSeen {
			return nil, fmt.Errorf("branch hierarchy output missing row for %s", childRef)
		}
		matrix[childRef] = candidates
	}
	return matrix, nil
}

func (self *BranchLoader) loadDirectBranchCandidates(branches []*models.Branch) (branchRelationMatrix, error) {
	matrix, err := self.loadBranchRelationMatrix(branches)
	if err != nil {
		return nil, err
	}
	return self.processLoadedBranchCandidates(matrix, branches)
}

func (self *BranchLoader) processLoadedBranchCandidates(
	matrix branchRelationMatrix,
	branches []*models.Branch,
) (branchRelationMatrix, error) {
	if !self.version.IsOlderThan(2, 41, 0) {
		return matrix, nil
	}

	directCandidates, err := self.loadLegacyDirectBranchCandidates(matrix, branches)
	if err != nil {
		return nil, err
	}
	if err := self.loadLegacyBranchDistances(directCandidates, branches); err != nil {
		return nil, err
	}
	return directCandidates, nil
}

func (self *BranchLoader) loadLegacyDirectBranchCandidates(
	matrix branchRelationMatrix,
	branches []*models.Branch,
) (branchRelationMatrix, error) {
	branchByRef := make(map[string]*models.Branch, len(branches))
	for _, branch := range branches {
		if branch.CommitHash != "" {
			branchByRef[branch.FullRefName()] = branch
		}
	}

	result := make(branchRelationMatrix, len(matrix))
	childRefs := slices.Sorted(maps.Keys(matrix))
	for _, childRef := range childRefs {
		candidatesByOID := map[string][]string{}
		for candidateRef := range matrix[childRef] {
			candidate := branchByRef[candidateRef]
			if candidate != nil {
				candidatesByOID[candidate.CommitHash] = append(candidatesByOID[candidate.CommitHash], candidateRef)
			}
		}

		candidateOIDs := slices.Sorted(maps.Keys(candidatesByOID))
		if len(candidateOIDs) <= 1 {
			result[childRef] = maps.Clone(matrix[childRef])
			continue
		}
		output, err := self.cmd.New(
			NewGitCmd("merge-base").
				Arg("--independent").
				Arg(candidateOIDs...).
				ToArgv(),
		).DontLog().RunWithOutput()
		if err != nil {
			return nil, err
		}

		direct := map[string]int{}
		for _, oid := range strings.Fields(output) {
			candidateRefs := candidatesByOID[oid]
			if len(candidateRefs) == 0 {
				return nil, fmt.Errorf("unexpected independent branch OID: %s", oid)
			}
			for _, candidateRef := range candidateRefs {
				direct[candidateRef] = matrix[childRef][candidateRef]
			}
		}
		if len(direct) == 0 {
			return nil, fmt.Errorf("branch hierarchy reduction returned no candidates for %s", childRef)
		}
		result[childRef] = direct
	}
	return result, nil
}

func (self *BranchLoader) loadLegacyBranchDistances(
	directCandidates branchRelationMatrix,
	branches []*models.Branch,
) error {
	branchByRef := make(map[string]*models.Branch, len(branches))
	for _, branch := range branches {
		if branch.CommitHash != "" {
			branchByRef[branch.FullRefName()] = branch
		}
	}

	childRefs := make([]string, 0, len(directCandidates))
	for childRef := range directCandidates {
		childRefs = append(childRefs, childRef)
	}
	slices.Sort(childRefs)
	for _, childRef := range childRefs {
		candidates := directCandidates[childRef]
		if len(candidates) <= 1 {
			continue
		}
		candidateRefs := make([]string, 0, len(candidates))
		for candidateRef := range candidates {
			candidateRefs = append(candidateRefs, candidateRef)
		}
		slices.Sort(candidateRefs)
		for _, candidateRef := range candidateRefs {
			child := branchByRef[childRef]
			candidate := branchByRef[candidateRef]
			if child == nil || candidate == nil {
				continue
			}
			output, err := self.cmd.New(
				NewGitCmd("rev-list").
					Arg("--count").
					Arg(candidate.CommitHash + ".." + child.CommitHash).
					ToArgv(),
			).DontLog().RunWithOutput()
			if err != nil {
				return err
			}
			distance, err := strconv.Atoi(strings.TrimSpace(output))
			if err != nil {
				return fmt.Errorf("parse branch distance for %s..%s: %w", candidateRef, childRef, err)
			}
			candidates[candidateRef] = distance
		}
	}
	return nil
}

func (self *BranchLoader) applyBranchHierarchy(branches []*models.Branch) []*models.Branch {
	for _, branch := range branches {
		branch.HierarchyDepth = 0
	}
	directCandidates, err := self.loadDirectBranchCandidates(branches)
	if err == nil {
		return inferBranchHierarchy(branches, directCandidates, self.UserConfig().Git.MainBranches)
	}

	self.Log.Errorf("failed to infer local branch hierarchy: %v", err)
	sortBranchesByNewest(branches)
	for index, branch := range branches {
		if branch.Head {
			return slices.Concat([]*models.Branch{branch}, branches[:index], branches[index+1:])
		}
	}
	return branches
}

func inferBranchHierarchy(
	branches []*models.Branch,
	directCandidates branchRelationMatrix,
	mainBranchNames []string,
) []*models.Branch {
	branchByRef := make(map[string]*models.Branch, len(branches))
	syntheticHeads := make([]*models.Branch, 0, 1)
	for _, branch := range branches {
		if branch.Head && branch.CommitHash == "" {
			syntheticHeads = append(syntheticHeads, branch)
			continue
		}
		branchByRef[branch.FullRefName()] = branch
	}

	candidatesByChild := make(branchRelationMatrix, len(directCandidates))
	unambiguousChildCount := map[string]int{}
	for childRef, candidates := range directCandidates {
		if branchByRef[childRef] == nil {
			continue
		}
		filtered := minimumDistanceCandidates(candidates, branchByRef)
		candidatesByChild[childRef] = filtered
		if len(filtered) == 1 {
			for parentRef := range filtered {
				unambiguousChildCount[parentRef]++
			}
		}
	}

	mainBranchSet := make(map[string]struct{}, len(mainBranchNames))
	for _, name := range mainBranchNames {
		mainBranchSet[strings.TrimPrefix(name, "refs/heads/")] = struct{}{}
	}

	parentByChild := make(map[string]string, len(candidatesByChild))
	childrenByParent := map[string][]*models.Branch{}
	for childRef, candidates := range candidatesByChild {
		if len(candidates) == 0 {
			continue
		}
		candidateRefs := make([]string, 0, len(candidates))
		for candidateRef := range candidates {
			candidateRefs = append(candidateRefs, candidateRef)
		}
		slices.SortFunc(candidateRefs, func(a, b string) int {
			if countComparison := unambiguousChildCount[b] - unambiguousChildCount[a]; countComparison != 0 {
				return countComparison
			}
			aIsMain := isMainBranch(branchByRef[a], mainBranchSet)
			bIsMain := isMainBranch(branchByRef[b], mainBranchSet)
			if aIsMain != bIsMain {
				if aIsMain {
					return -1
				}
				return 1
			}
			if timestampComparison := compareTimestampsAscending(branchByRef[a], branchByRef[b]); timestampComparison != 0 {
				return timestampComparison
			}
			return strings.Compare(branchByRef[a].Name, branchByRef[b].Name)
		})
		parentRef := candidateRefs[0]
		parentByChild[childRef] = parentRef
		childrenByParent[parentRef] = append(childrenByParent[parentRef], branchByRef[childRef])
	}

	for parentRef := range childrenByParent {
		sortBranchesByNewest(childrenByParent[parentRef])
	}

	roots := make([]*models.Branch, 0, len(branchByRef))
	for ref, branch := range branchByRef {
		if _, hasParent := parentByChild[ref]; !hasParent {
			roots = append(roots, branch)
		}
	}
	sortBranchesByNewest(roots)
	promoteHeadRoot(roots, parentByChild, branchByRef)

	result := make([]*models.Branch, 0, len(branches))
	for _, branch := range syntheticHeads {
		branch.HierarchyDepth = 0
		result = append(result, branch)
	}
	var appendSubtree func(*models.Branch, int)
	appendSubtree = func(branch *models.Branch, depth int) {
		branch.HierarchyDepth = depth
		result = append(result, branch)
		for _, child := range childrenByParent[branch.FullRefName()] {
			appendSubtree(child, depth+1)
		}
	}
	for _, root := range roots {
		appendSubtree(root, 0)
	}
	return result
}

func minimumDistanceCandidates(
	candidates map[string]int,
	branchByRef map[string]*models.Branch,
) map[string]int {
	result := make(map[string]int, len(candidates))
	for candidateRef, distance := range candidates {
		if branchByRef[candidateRef] != nil {
			result[candidateRef] = distance
		}
	}
	if len(result) <= 1 {
		return result
	}

	minimumDistance := -1
	for _, distance := range result {
		if distance >= 0 && (minimumDistance == -1 || distance < minimumDistance) {
			minimumDistance = distance
		}
	}
	if minimumDistance == -1 {
		return result
	}
	for candidateRef, distance := range result {
		if distance != minimumDistance {
			delete(result, candidateRef)
		}
	}
	return result
}

func isMainBranch(branch *models.Branch, mainBranchSet map[string]struct{}) bool {
	if branch == nil {
		return false
	}
	_, ok := mainBranchSet[branch.Name]
	return ok
}

func compareTimestampsAscending(a, b *models.Branch) int {
	if a.CommitUnixTimestamp < b.CommitUnixTimestamp {
		return -1
	}
	if a.CommitUnixTimestamp > b.CommitUnixTimestamp {
		return 1
	}
	return 0
}

func sortBranchesByNewest(branches []*models.Branch) {
	slices.SortFunc(branches, func(a, b *models.Branch) int {
		if a.CommitUnixTimestamp > b.CommitUnixTimestamp {
			return -1
		}
		if a.CommitUnixTimestamp < b.CommitUnixTimestamp {
			return 1
		}
		return strings.Compare(a.Name, b.Name)
	})
}

func promoteHeadRoot(
	roots []*models.Branch,
	parentByChild map[string]string,
	branchByRef map[string]*models.Branch,
) {
	for _, branch := range branchByRef {
		if !branch.Head {
			continue
		}
		rootRef := branch.FullRefName()
		for parentRef := parentByChild[rootRef]; parentRef != ""; parentRef = parentByChild[rootRef] {
			rootRef = parentRef
		}
		moveBranchToFront(roots, rootRef)
		return
	}
}

func moveBranchToFront(branches []*models.Branch, ref string) {
	for index, branch := range branches {
		if branch.FullRefName() == ref {
			copy(branches[1:index+1], branches[:index])
			branches[0] = branch
			return
		}
	}
}
