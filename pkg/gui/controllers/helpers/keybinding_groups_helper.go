package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type KeybindingGroupsHelper struct {
	c *HelperCommon
}

func NewKeybindingGroupsHelper(c *HelperCommon) *KeybindingGroupsHelper {
	return &KeybindingGroupsHelper{c: c}
}

// ContextByName resolves a switchTo enum value to its types.Context, or nil
// if unknown. Validation in pkg/config rejects unknown names at load time;
// the nil branch protects against runtime drift.
func (self *KeybindingGroupsHelper) ContextByName(name string) types.Context {
	tree := self.c.Contexts()
	switch name {
	case "status":
		return tree.Status
	case "files":
		return tree.Files
	case "localBranches":
		return tree.Branches
	case "remotes":
		return tree.Remotes
	case "tags":
		return tree.Tags
	case "localCommits":
		return tree.LocalCommits
	case "subCommits":
		return tree.SubCommits
	case "reflog":
		return tree.ReflogCommits
	case "stash":
		return tree.Stash
	case "submodules":
		return tree.Submodules
	case "worktrees":
		return tree.Worktrees
	case "commitFiles":
		return tree.CommitFiles
	}
	return nil
}
