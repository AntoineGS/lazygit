package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// ResetToRefTitle returns a dynamic chord-menu title closure for
// "Reset to <ref>" groups. When getRef returns (_, false) — typically
// because no ref is selected, the model is empty, or git is unavailable —
// the closure falls back to the supplied fallback string.
func ResetToRefTitle(c *HelperCommon, fallback string, getRef func() (string, bool)) func() string {
	return func() string {
		ref, ok := getRef()
		if !ok {
			return fallback
		}
		return c.Tr.ResetTo + " " + ref
	}
}

// RebasingTitle returns a dynamic chord-menu title closure for the
// "Rebase '<branch>'" group. If withMarkedBase is true, the closure
// switches between c.Tr.RebasingTitle and c.Tr.RebasingFromBaseCommitTitle
// based on whether a base commit is currently marked. The static
// c.Tr.RebasingTitle is returned if git is unavailable or no branches
// are loaded yet.
func RebasingTitle(c *HelperCommon, withMarkedBase bool) func() string {
	return func() string {
		if c.Git() == nil || len(c.Model().Branches) == 0 {
			return c.Tr.RebasingTitle
		}
		baseTitle := c.Tr.RebasingTitle
		if withMarkedBase && c.Modes().MarkedBaseCommit.GetHash() != "" {
			baseTitle = c.Tr.RebasingFromBaseCommitTitle
		}
		return utils.ResolvePlaceholderString(
			baseTitle,
			map[string]string{"checkedOutBranch": c.Model().Branches[0].Name},
		)
	}
}
