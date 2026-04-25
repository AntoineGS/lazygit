package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemoteBranchesController struct {
	baseController
	*ListControllerTrait[*models.RemoteBranch]
	c *ControllerCommon
}

var _ types.IController = &RemoteBranchesController{}

func NewRemoteBranchesController(
	c *ControllerCommon,
) *RemoteBranchesController {
	return &RemoteBranchesController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().RemoteBranches,
			c.Contexts().RemoteBranches.GetSelected,
			c.Contexts().RemoteBranches.GetSelectedItems,
		),
		c: c,
	}
}

func (self *RemoteBranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.checkoutBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.RemoteBranchCheckoutTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.newLocalBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeRegular),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeRegular)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Merge,
			Tooltip:           self.c.Tr.MergeBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeNonFFwd),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeNonFastForward)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Non-fast-forward merge",
			Tooltip:           self.c.Tr.RegularMergeNonFastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeFastForward),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeFastForward)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Fast-forward only merge",
			Tooltip:           self.c.Tr.RegularMergeFastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeSquash),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeSquash)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Squash merge (uncommitted)",
			Tooltip:           self.c.Tr.SquashMergeUncommittedTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeSquashCommitted),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeSquashCommitted)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Squash merge (committed)",
			Tooltip:           self.c.Tr.SquashMergeCommittedTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranchSimple),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebaseSimple)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RebaseBranch,
			Tooltip:           self.c.Tr.RebaseBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranchInteractive),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebaseInteractive)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Interactive rebase",
			Tooltip:           self.c.Tr.InteractiveRebaseTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.RebaseBranchOntoBase),
			Handler:     opts.Guards.OutsideFilterMode(self.rebaseOntoBaseBranch),
			Description: "Rebase onto base branch",
			Tooltip:     self.c.Tr.RebaseOntoBaseBranchTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.delete),
			GetDisabledReason: self.require(self.itemRangeSelected()),
			Description:       self.c.Tr.Delete,
			Tooltip:           self.c.Tr.DeleteRemoteBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:           self.withItem(self.setAsUpstream),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.SetAsUpstream,
			Tooltip:           self.c.Tr.SetAsUpstreamTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SortOrder),
			Handler:     self.createSortMenu,
			Description: self.c.Tr.SortOrder,
			OpensMenu:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MixedResetToRef),
			Handler:           self.withItem(self.gitMixedResetToRef),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Mixed reset",
			Tooltip:           self.c.Tr.ResetMixedTooltip,
			ChordPopupExtra:   self.gitResetPreview(style.FgRed, "mixed"),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.SoftResetToRef),
			Handler:           self.withItem(self.gitSoftResetToRef),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.SoftReset,
			Tooltip:           self.c.Tr.ResetSoftTooltip,
			ChordPopupExtra:   self.gitResetPreview(style.FgRed, "soft"),
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.HardResetToRef),
			Handler:           self.withItem(self.gitHardResetToRef),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.HardReset,
			Tooltip:           self.c.Tr.ResetHardTooltip,
			ChordPopupExtra:   self.gitResetPreview(style.FgRed, "hard"),
		},
		{
			Key: opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler: self.withItem(func(selectedBranch *models.RemoteBranch) error {
				return self.c.Helpers().Diff.OpenDiffToolForRef(selectedBranch)
			}),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
	}
}

func (self *RemoteBranchesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			remoteBranch := self.context().GetSelected()
			if remoteBranch == nil {
				task = types.NewRenderStringTask("No branches for this remote")
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(remoteBranch.FullRefName())
				task = types.NewRunCommandTask(cmdObj.GetCmd())
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Remote Branch",
					Task:  task,
				},
			})
		})
	}
}

func (self *RemoteBranchesController) context() *context.RemoteBranchesContext {
	return self.c.Contexts().RemoteBranches
}

func (self *RemoteBranchesController) delete(selectedBranches []*models.RemoteBranch) error {
	return self.c.Helpers().BranchesHelper.ConfirmDeleteRemote(selectedBranches, true)
}

func (self *RemoteBranchesController) mergeRegular(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(selectedBranch.FullName(), git_commands.MERGE_VARIANT_REGULAR)
}

func (self *RemoteBranchesController) mergeNonFastForward(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(selectedBranch.FullName(), git_commands.MERGE_VARIANT_NON_FAST_FORWARD)
}

func (self *RemoteBranchesController) mergeFastForward(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(selectedBranch.FullName(), git_commands.MERGE_VARIANT_FAST_FORWARD)
}

func (self *RemoteBranchesController) mergeSquashCommitted(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformSquashMergeCommitted(selectedBranch.FullName())
}

func (self *RemoteBranchesController) mergeSquash(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformSquashMerge(selectedBranch.FullName())
}

func (self *RemoteBranchesController) rebaseSimple(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(selectedBranch.FullName(), helpers.RebaseVariantSimple)
}

func (self *RemoteBranchesController) rebaseInteractive(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(selectedBranch.FullName(), helpers.RebaseVariantInteractive)
}

func (self *RemoteBranchesController) rebaseOntoBaseBranch() error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef("", helpers.RebaseVariantOntoBase)
}

func (self *RemoteBranchesController) createSortMenu() error {
	return self.c.Helpers().Refs.CreateSortOrderMenu(
		[]string{"alphabetical", "date"},
		self.c.Tr.SortOrderPromptRemoteBranches,
		func(sortOrder string) error {
			if self.c.UserConfig().Git.RemoteBranchSortOrder != sortOrder {
				self.c.UserConfig().Git.RemoteBranchSortOrder = sortOrder
				self.c.Contexts().RemoteBranches.SetSelection(0)
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.REMOTES}})
			}
			return nil
		},
		self.c.UserConfig().Git.RemoteBranchSortOrder)
}

func (self *RemoteBranchesController) gitMixedResetToRef(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.FullName(), selectedBranch.FullRefName(), "mixed")
}

func (self *RemoteBranchesController) gitSoftResetToRef(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.FullName(), selectedBranch.FullRefName(), "soft")
}

func (self *RemoteBranchesController) gitHardResetToRef(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.FullName(), selectedBranch.FullRefName(), "hard")
}

func (self *RemoteBranchesController) gitResetPreview(s style.TextStyle, strength string) string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil {
		return ""
	}
	return s.Sprintf("reset --%s %s", strength, branch.FullName())
}

func (self *RemoteBranchesController) setAsUpstream(selectedBranch *models.RemoteBranch) error {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()

	message := utils.ResolvePlaceholderString(
		self.c.Tr.SetUpstreamMessage,
		map[string]string{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.FullName(),
		},
	)

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.SetUpstreamTitle,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.SetBranchUpstream)
			if err := self.c.Git().Branch.SetUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return err
			}

			self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			return nil
		},
	})

	return nil
}

func (self *RemoteBranchesController) newLocalBranch(selectedBranch *models.RemoteBranch) error {
	// will set to the remote's branch name without the remote name
	nameSuggestion := strings.SplitAfterN(selectedBranch.RefName(), "/", 2)[1]

	return self.c.Helpers().Refs.NewBranch(selectedBranch.RefName(), selectedBranch.RefName(), nameSuggestion)
}

func (self *RemoteBranchesController) checkoutBranch(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.CheckoutRemoteBranch(selectedBranch.FullName(), selectedBranch.Name)
}
