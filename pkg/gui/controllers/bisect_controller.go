package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type BisectController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon
}

var _ types.IController = &BisectController{}

func NewBisectController(
	c *ControllerCommon,
) *BisectController {
	return &BisectController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().LocalCommits,
			c.Contexts().LocalCommits.GetSelected,
			c.Contexts().LocalCommits.GetSelectedItems,
		),
	}
}

func (self *BisectController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	midDisabled := func() *types.DisabledReason {
		if self.c.Git() == nil || !self.c.Git().Bisect.GetInfo().Started() {
			return &types.DisabledReason{AllowFurtherDispatching: true}
		}
		return nil
	}
	startDisabled := func() *types.DisabledReason {
		if self.c.Git() != nil && self.c.Git().Bisect.GetInfo().Started() {
			return &types.DisabledReason{AllowFurtherDispatching: true}
		}
		return nil
	}

	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectMarkBad),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectMarkBad)),
			Description:       self.c.Tr.BisectMarkBad,
			DescriptionFunc:   self.bisectMarkBadLabel,
			GetDisabledReason: self.bisectMarkDisabledReason(midDisabled),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectMarkGood),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectMarkGood)),
			Description:       self.c.Tr.BisectMarkGood,
			DescriptionFunc:   self.bisectMarkGoodLabel,
			GetDisabledReason: self.bisectMarkDisabledReason(midDisabled),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectSkipCurrent),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectSkipCurrent)),
			Description:       self.c.Tr.BisectSkipCurrent,
			DescriptionFunc:   self.bisectSkipCurrentLabel,
			GetDisabledReason: self.bisectMarkDisabledReason(midDisabled),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectSkipSelected),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectSkipSelected)),
			Description:       self.c.Tr.BisectSkipSelected,
			DescriptionFunc:   self.bisectSkipSelectedLabel,
			GetDisabledReason: self.bisectSkipSelectedDisabledReason,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectReset),
			Handler:           opts.Guards.OutsideFilterMode(self.bisectReset),
			Description:       self.c.Tr.Bisect.ResetOption,
			GetDisabledReason: midDisabled,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectStartMarkBad),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectStartMarkBad)),
			Description:       self.c.Tr.BisectStartMarkBad,
			DescriptionFunc:   self.bisectStartMarkBadLabel,
			GetDisabledReason: self.bisectStartMarkDisabledReason(startDisabled),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectStartMarkGood),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.bisectStartMarkGood)),
			Description:       self.c.Tr.BisectStartMarkGood,
			DescriptionFunc:   self.bisectStartMarkGoodLabel,
			GetDisabledReason: self.bisectStartMarkDisabledReason(startDisabled),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.BisectChooseTerms),
			Handler:           opts.Guards.OutsideFilterMode(self.bisectChooseTerms),
			Description:       self.c.Tr.Bisect.ChooseTerms,
			GetDisabledReason: startDisabled,
		},
	}
}

// modeReason returns AllowFurtherDispatching for mode mismatches; the
// secondary single-item reason does NOT, so a wrong-selection produces
// a toast rather than falling through.
func (self *BisectController) bisectMarkDisabledReason(modeReason func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if r := modeReason(); r != nil {
			return r
		}
		info := self.c.Git().Bisect.GetInfo()
		bisecting := info.GetCurrentHash() != ""
		if !bisecting {
			return self.require(self.singleItemSelected())()
		}
		return nil
	}
}

func (self *BisectController) bisectStartMarkDisabledReason(modeReason func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if r := modeReason(); r != nil {
			return r
		}
		return self.require(self.singleItemSelected())()
	}
}

func (self *BisectController) bisectSkipSelectedDisabledReason() *types.DisabledReason {
	if self.c.Git() == nil || !self.c.Git().Bisect.GetInfo().Started() {
		return &types.DisabledReason{AllowFurtherDispatching: true}
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil || info.GetCurrentHash() == "" || info.GetCurrentHash() == commit.Hash() {
		return &types.DisabledReason{}
	}
	return self.require(self.singleItemSelected())()
}

func (self *BisectController) bisectMarkContext(info *git_commands.BisectInfo, commit *models.Commit) (selectCurrentAfter, waitToReselect bool, hashToMark string) {
	// If there's no 'current' bisect commit, or the selected commit IS
	// the current one, we need to jump to the next 'current' commit
	// after the action.
	selectCurrentAfter = info.GetCurrentHash() == "" || info.GetCurrentHash() == commit.Hash()
	// Wait to reselect if our bisect commits aren't ancestors of the
	// 'start' ref, because we'll be reloading commits in that case.
	waitToReselect = selectCurrentAfter && !self.c.Git().Bisect.ReachableFromStart(info)
	bisecting := info.GetCurrentHash() != ""
	hashToMark = lo.Ternary(bisecting, info.GetCurrentHash(), commit.Hash())
	return selectCurrentAfter, waitToReselect, hashToMark
}

func (self *BisectController) bisectMarkBad(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	selectCurrentAfter, waitToReselect, hashToMark := self.bisectMarkContext(info, commit)
	self.c.LogAction(self.c.Tr.Actions.BisectMark)
	if err := self.c.Git().Bisect.Mark(hashToMark, info.NewTerm()); err != nil {
		return err
	}
	return self.afterMark(selectCurrentAfter, waitToReselect)
}

func (self *BisectController) bisectMarkGood(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	selectCurrentAfter, waitToReselect, hashToMark := self.bisectMarkContext(info, commit)
	self.c.LogAction(self.c.Tr.Actions.BisectMark)
	if err := self.c.Git().Bisect.Mark(hashToMark, info.OldTerm()); err != nil {
		return err
	}
	return self.afterMark(selectCurrentAfter, waitToReselect)
}

func (self *BisectController) bisectSkipCurrent(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	selectCurrentAfter, waitToReselect, hashToMark := self.bisectMarkContext(info, commit)
	self.c.LogAction(self.c.Tr.Actions.BisectSkip)
	if err := self.c.Git().Bisect.Skip(hashToMark); err != nil {
		return err
	}
	return self.afterMark(selectCurrentAfter, waitToReselect)
}

func (self *BisectController) bisectSkipSelected(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	selectCurrentAfter, waitToReselect, _ := self.bisectMarkContext(info, commit)
	self.c.LogAction(self.c.Tr.Actions.BisectSkip)
	if err := self.c.Git().Bisect.Skip(commit.Hash()); err != nil {
		return err
	}
	return self.afterMark(selectCurrentAfter, waitToReselect)
}

func (self *BisectController) bisectReset() error {
	return self.c.Helpers().Bisect.Reset()
}

func (self *BisectController) bisectStartMarkBad(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	self.c.LogAction(self.c.Tr.Actions.StartBisect)
	if err := self.c.Git().Bisect.Start(); err != nil {
		return err
	}
	if err := self.c.Git().Bisect.Mark(commit.Hash(), info.NewTerm()); err != nil {
		return err
	}
	self.c.Helpers().Bisect.PostBisectCommandRefresh()
	return nil
}

func (self *BisectController) bisectStartMarkGood(commit *models.Commit) error {
	info := self.c.Git().Bisect.GetInfo()
	self.c.LogAction(self.c.Tr.Actions.StartBisect)
	if err := self.c.Git().Bisect.Start(); err != nil {
		return err
	}
	if err := self.c.Git().Bisect.Mark(commit.Hash(), info.OldTerm()); err != nil {
		return err
	}
	self.c.Helpers().Bisect.PostBisectCommandRefresh()
	return nil
}

func (self *BisectController) bisectChooseTerms() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.Bisect.OldTermPrompt,
		HandleConfirm: func(oldTerm string) error {
			self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.Bisect.NewTermPrompt,
				HandleConfirm: func(newTerm string) error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.StartWithTerms(oldTerm, newTerm); err != nil {
						return err
					}
					self.c.Helpers().Bisect.PostBisectCommandRefresh()
					return nil
				},
			})
			return nil
		},
	})
	return nil
}

func (self *BisectController) bisectMarkBadLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectMarkBad
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectMarkBad
	}
	bisecting := info.GetCurrentHash() != ""
	hashToMark := lo.Ternary(bisecting, info.GetCurrentHash(), commit.Hash())
	return fmt.Sprintf(self.c.Tr.Bisect.Mark, utils.ShortHash(hashToMark), info.NewTerm())
}

func (self *BisectController) bisectMarkGoodLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectMarkGood
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectMarkGood
	}
	bisecting := info.GetCurrentHash() != ""
	hashToMark := lo.Ternary(bisecting, info.GetCurrentHash(), commit.Hash())
	return fmt.Sprintf(self.c.Tr.Bisect.Mark, utils.ShortHash(hashToMark), info.OldTerm())
}

func (self *BisectController) bisectSkipCurrentLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectSkipCurrent
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectSkipCurrent
	}
	bisecting := info.GetCurrentHash() != ""
	hashToMark := lo.Ternary(bisecting, info.GetCurrentHash(), commit.Hash())
	return fmt.Sprintf(self.c.Tr.Bisect.SkipCurrent, utils.ShortHash(hashToMark))
}

func (self *BisectController) bisectSkipSelectedLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectSkipSelected
	}
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectSkipSelected
	}
	return fmt.Sprintf(self.c.Tr.Bisect.SkipSelected, commit.ShortHash())
}

func (self *BisectController) bisectStartMarkBadLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectStartMarkBad
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectStartMarkBad
	}
	return fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortHash(), info.NewTerm())
}

func (self *BisectController) bisectStartMarkGoodLabel() string {
	if self.c.Git() == nil {
		return self.c.Tr.BisectStartMarkGood
	}
	info := self.c.Git().Bisect.GetInfo()
	commit := self.context().GetSelected()
	if commit == nil {
		return self.c.Tr.BisectStartMarkGood
	}
	return fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortHash(), info.OldTerm())
}

func (self *BisectController) showBisectCompleteMessage(candidateHashes []string) error {
	prompt := self.c.Tr.Bisect.CompletePrompt
	if len(candidateHashes) > 1 {
		prompt = self.c.Tr.Bisect.CompletePromptIndeterminate
	}

	formattedCommits, err := self.c.Git().Commit.GetCommitsOneline(candidateHashes)
	if err != nil {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Bisect.CompleteTitle,
		Prompt: fmt.Sprintf(prompt, strings.TrimSpace(formattedCommits)),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.c.Git().Bisect.Reset(); err != nil {
				return err
			}

			self.c.Helpers().Bisect.PostBisectCommandRefresh()
			return nil
		},
	})

	return nil
}

func (self *BisectController) afterMark(selectCurrent bool, waitToReselect bool) error {
	done, candidateHashes, err := self.c.Git().Bisect.IsDone()
	if err != nil {
		return err
	}

	if err := self.afterBisectMarkRefresh(selectCurrent, waitToReselect); err != nil {
		return err
	}

	if done {
		return self.showBisectCompleteMessage(candidateHashes)
	}

	return nil
}

func (self *BisectController) afterBisectMarkRefresh(selectCurrent bool, waitToReselect bool) error {
	selectFn := func() {
		if selectCurrent {
			self.selectCurrentBisectCommit()
		}
	}

	if waitToReselect {
		self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{}, Then: selectFn})
		return nil
	}

	selectFn()

	self.c.Helpers().Bisect.PostBisectCommandRefresh()
	return nil
}

func (self *BisectController) selectCurrentBisectCommit() {
	info := self.c.Git().Bisect.GetInfo()
	if info.GetCurrentHash() != "" {
		// find index of commit with that hash, move cursor to that.
		for i, commit := range self.c.Model().Commits {
			if commit.Hash() == info.GetCurrentHash() {
				self.context().SetSelection(i)
				self.context().HandleFocus(types.OnFocusOpts{})
				break
			}
		}
	}
}

func (self *BisectController) context() *context.LocalCommitsContext {
	return self.c.Contexts().LocalCommits
}
