package git_commands

import (
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// Holds parsed values from a single %(ahead-behind:<base>) field.
type aheadBehind struct {
	baseRevision string
	ahead        int
	behind       int
}

type branchAheadBehind struct {
	refName      string
	aheadBehinds []aheadBehind
}

// Parses output produced by:
//
//	git for-each-ref --format='%(refname)\x00%(ahead-behind:<base1>)\x00...' refs/heads
//
// Lines whose NUL-split column count doesn't match (1 + len(baseRevisions)) are dropped.
// Blank lines are ignored.
// Individual malformed ahead-behind fields are omitted.
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

func parseAheadBehindField(s string) (aheadBehind, bool) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return aheadBehind{}, false
	}
	ahead, err1 := strconv.Atoi(parts[0])
	behind, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return aheadBehind{}, false
	}
	return aheadBehind{ahead: ahead, behind: behind}, true
}

// Picks the "closest" base by smallest ahead value (commits the branch
// has that the base doesn't = roughly "since fork point") and returns
// its behind value.
// Ties are broken by index order
func selectBehindForBranch(values []aheadBehind) int {
	if len(values) == 0 {
		return 0
	}
	return lo.MinBy(values, func(a, b aheadBehind) bool {
		return a.ahead < b.ahead
	}).behind
}

// The output format is:
//
//	<refname>\x00<ahead> <behind>\x00<ahead> <behind>...\n
//
// with one ahead-behind field per base, in the same order as mainBranchRefs.
//
// Requires git >= 2.41 (when %(ahead-behind:...) was added).
func buildAheadBehindForEachRefArgs(mainBranchRefs []string) []string {
	formatParts := make([]string, 0, 1+len(mainBranchRefs))
	formatParts = append(formatParts, "%(refname)")
	for _, ref := range mainBranchRefs {
		formatParts = append(formatParts, "%(ahead-behind:"+ref+")")
	}
	format := strings.Join(formatParts, "%00")

	return NewGitCmd("for-each-ref").
		Arg("--format=" + format).
		Arg("refs/heads").
		ToArgv()
}
