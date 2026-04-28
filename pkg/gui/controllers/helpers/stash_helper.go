package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StashHelper struct {
	c *HelperCommon
}

func NewStashHelper(c *HelperCommon) *StashHelper {
	return &StashHelper{c: c}
}

// PopStash pops the given stash entry. Behavior identical to the old
// StashController.handleStashPop body — lifted into a helper so callers
// outside StashController can invoke it.
func (self *StashHelper) PopStash(stashEntry *models.StashEntry) error {
	pop := func() error {
		self.c.LogAction(self.c.Tr.Actions.PopStash)
		self.c.LogCommand(fmt.Sprintf(self.c.Tr.Log.PoppingStash, stashEntry.Hash), false)
		err := self.c.Git().Stash.Pop(stashEntry.Index)
		self.postStashRefresh()
		if err != nil {
			return err
		}
		if self.c.UserConfig().Gui.SwitchToFilesAfterStashPop {
			self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
		}
		return nil
	}

	if self.c.UserConfig().Gui.SkipStashWarning {
		return pop()
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashPop,
		Prompt: self.c.Tr.SurePopStashEntry,
		HandleConfirm: func() error {
			return pop()
		},
	})

	return nil
}

// PopStashWithDefault is the universal-binding entry point. Uses the
// stash entry currently selected in the stash view if one exists;
// otherwise falls back to the topmost stash, mirroring `git stash pop`
// with no argument.
func (self *StashHelper) PopStashWithDefault() error {
	stashEntries := self.c.Model().StashEntries
	if len(stashEntries) == 0 {
		return nil
	}

	target := self.c.Contexts().Stash.GetSelected()
	if target == nil {
		target = stashEntries[0]
	}

	return self.PopStash(target)
}

// StashAllChanges stashes every working-tree change. Behavior identical
// to the old FilesController.stash body (which delegated to
// handleStashSave with Push + StashAllChanges) — inlined here so
// non-files-view callers can invoke it without the controller.
func (self *StashHelper) StashAllChanges() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			self.c.LogAction(self.c.Tr.Actions.StashAllChanges)

			if err := self.c.Git().Stash.Push(stashComment); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
			return nil
		},
		AllowEmptyInput: true,
	})

	return nil
}

// postStashRefresh is the same private helper StashController used to
// call. Lifted verbatim.
func (self *StashHelper) postStashRefresh() {
	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
}
