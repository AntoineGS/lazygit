package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type BranchesController struct {
	baseController
	*ListControllerTrait[*models.Branch]
	c *ControllerCommon
}

var _ types.IController = &BranchesController{}

func NewBranchesController(
	c *ControllerCommon,
) *BranchesController {
	return &BranchesController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().Branches,
			c.Contexts().Branches.GetSelected,
			c.Contexts().Branches.GetSelectedItems,
		),
	}
}

func (self *BranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.Select),
			Handler: self.withItem(self.press),
			GetDisabledReason: self.require(
				self.singleItemSelected(),
				self.notPulling,
			),
			Description:     self.c.Tr.Checkout,
			Tooltip:         self.c.Tr.CheckoutTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.newBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MoveCommitsToNewBranch),
			Handler:           self.c.Helpers().Refs.MoveCommitsToNewBranch,
			GetDisabledReason: self.c.Helpers().Refs.CanMoveCommitsToNewBranch,
			Description:       self.c.Tr.MoveCommitsToNewBranch,
			Tooltip:           self.c.Tr.MoveCommitsToNewBranchTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CreatePullRequest),
			Handler:           self.withItem(self.handleCreatePullRequest),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreatePullRequest,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.ViewPullRequestOptions),
			Handler:           self.withItem(self.handleCreatePullRequestMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreatePullRequestOptions,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.OpenPullRequestInBrowser),
			Handler:           self.withItem(self.openPRInBrowser),
			GetDisabledReason: self.require(self.singleItemSelected(self.branchHasPR)),
			Description:       self.c.Tr.OpenPullRequestInBrowser,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CopyPullRequestURL),
			Handler:           self.copyPullRequestURL,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CopyPullRequestURL,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CheckoutBranchByName),
			Handler:     self.checkoutByName,
			Description: self.c.Tr.CheckoutByName,
			Tooltip:     self.c.Tr.CheckoutByNameTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CheckoutPreviousBranch),
			Handler:     self.checkoutPreviousBranch,
			Description: self.c.Tr.CheckoutPreviousBranch,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.ForceCheckoutBranch),
			Handler:           self.forceCheckout,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ForceCheckout,
			Tooltip:           self.c.Tr.ForceCheckoutTooltip,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.DeleteLocalBranch),
			Handler: self.withItems(self.localDelete),
			GetDisabledReason: self.require(
				self.itemRangeSelected(self.branchesAreReal),
				self.itemRangeSelected(self.notDeletingCheckedOutBranch),
			),
			Description:     self.c.Tr.DeleteLocalBranch,
			DescriptionFunc: self.deleteBranchDescriptionFunc(self.c.Tr.DeleteLocalBranch, self.c.Tr.DeleteLocalBranches),
			Tooltip:         self.c.Tr.BranchDeleteTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.DeleteRemoteBranch),
			Handler: self.withItems(self.remoteDelete),
			GetDisabledReason: self.require(
				self.itemRangeSelected(self.branchesAreReal),
				self.itemRangeSelected(self.allBranchesHaveUpstream),
			),
			Description:     self.c.Tr.DeleteRemoteBranch,
			DescriptionFunc: self.deleteBranchDescriptionFunc(self.c.Tr.DeleteRemoteBranch, self.c.Tr.DeleteRemoteBranches),
			Tooltip:         self.c.Tr.DeleteRemoteBranchTooltip,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.DeleteLocalAndRemoteBranch),
			Handler: self.withItems(self.localAndRemoteDelete),
			GetDisabledReason: self.require(
				self.itemRangeSelected(self.branchesAreReal),
				self.itemRangeSelected(self.notDeletingCheckedOutBranch),
				self.itemRangeSelected(self.allBranchesHaveUpstream),
			),
			Description:     self.c.Tr.DeleteLocalAndRemoteBranch,
			DescriptionFunc: self.deleteBranchDescriptionFunc(self.c.Tr.DeleteLocalAndRemoteBranch, self.c.Tr.DeleteLocalAndRemoteBranches),
			Tooltip:         self.c.Tr.BranchDeleteTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranchSimple),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebaseSimple)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notRebasingOntoSelf)),
			Description:       "Simple rebase",
			DescriptionFunc:   self.rebaseBranchSimpleDescriptionFunc(),
			Tooltip:           self.c.Tr.RebaseBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranchInteractive),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebaseInteractive)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notRebasingOntoSelf)),
			Description:       "Interactive rebase",
			DescriptionFunc:   self.rebaseBranchInteractiveDescriptionFunc(),
			Tooltip:           self.c.Tr.InteractiveRebaseTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranchOntoBase),
			Handler:           opts.Guards.OutsideFilterMode(self.rebaseOntoBaseBranch),
			GetDisabledReason: self.require(self.singleItemSelected(self.canRebaseOntoBase)),
			Description:       "Rebase onto base branch",
			DescriptionFunc:   self.rebaseOntoBaseDescriptionFunc(),
			Tooltip:           self.c.Tr.RebaseOntoBaseBranchTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeRegular),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeRegular)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself)),
			Description:       self.c.Tr.Merge,
			Tooltip:           self.c.Tr.MergeBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeNonFFwd),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeNonFastForward)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself, self.nonFastForwardMergeApplicable)),
			Description:       self.c.Tr.RegularMergeNonFastForward,
			Tooltip:           self.c.Tr.RegularMergeNonFastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeFastForward),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeFastForward)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself, self.fastForwardOnlyMergeApplicable)),
			Description:       self.c.Tr.RegularMergeFastForward,
			Tooltip:           self.c.Tr.RegularMergeFastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeSquash),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeSquash)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself)),
			Description:       "Squash merge (uncommitted)",
			Tooltip:           self.c.Tr.SquashMergeUncommittedTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeSquashCommitted),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.mergeSquashCommitted)),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself)),
			Description:       "Squash merge (committed)",
			Tooltip:           self.c.Tr.SquashMergeCommittedTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.FastForward),
			Handler:           self.withItem(self.fastForward),
			GetDisabledReason: self.require(self.singleItemSelected(self.branchIsReal)),
			Description:       self.c.Tr.FastForward,
			Tooltip:           self.c.Tr.FastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CreateTag),
			Handler:           self.withItem(self.createTag),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewTag,
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
			Key:               opts.GetKey(opts.Config.Branches.RenameBranch),
			Handler:           self.withItem(self.rename),
			GetDisabledReason: self.require(self.singleItemSelected(self.branchIsReal)),
			Description:       self.c.Tr.RenameBranch,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.ViewDivergenceFromUpstream),
			Handler: self.withItem(self.viewDivergenceFromUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     self.c.Tr.ViewDivergenceFromUpstream,
			DescriptionFunc: self.viewDivergenceFromUpstreamDescriptionFunc(),
			Tooltip:         self.c.Tr.ViewBranchUpstreamOptionsTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.ViewDivergenceFromBase),
			Handler:           self.withItem(self.viewDivergenceFromBaseBranch),
			GetDisabledReason: self.require(self.singleItemSelected(self.canViewDivergenceFromBase)),
			Description:       "View divergence from base branch",
			DescriptionFunc:   self.viewDivergenceFromBaseDescriptionFunc(),
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:           self.withItem(self.setUpstream),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.SetUpstream,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.UnsetUpstream),
			Handler: self.withItem(self.unsetUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.branchIsTrackingRemote),
			),
			Description: self.c.Tr.UnsetUpstream,
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.ResetUpstreamMixed),
			Handler: self.withItem(self.gitMixedResetToUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     "Mixed reset to upstream",
			DescriptionFunc: self.upstreamResetDescriptionFunc("Mixed reset to upstream"),
			Tooltip:         self.c.Tr.ResetMixedTooltip,
			ChordPopupExtra: self.upstreamResetPreview("mixed"),
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.ResetUpstreamSoft),
			Handler: self.withItem(self.gitSoftResetToUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     "Soft reset to upstream",
			DescriptionFunc: self.upstreamResetDescriptionFunc("Soft reset to upstream"),
			Tooltip:         self.c.Tr.ResetSoftTooltip,
			ChordPopupExtra: self.upstreamResetPreview("soft"),
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.ResetUpstreamHard),
			Handler: self.withItem(self.gitHardResetToUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     "Hard reset to upstream",
			DescriptionFunc: self.upstreamResetDescriptionFunc("Hard reset to upstream"),
			Tooltip:         self.c.Tr.ResetHardTooltip,
			ChordPopupExtra: self.upstreamResetPreview("hard"),
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.RebaseUpstreamSimple),
			Handler: self.withItem(self.rebaseSimpleOntoUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     "Simple rebase onto upstream",
			DescriptionFunc: self.upstreamRebaseDescriptionFunc("Simple rebase onto upstream"),
		},
		{
			Key:     opts.GetKey(opts.Config.Branches.RebaseUpstreamInteractive),
			Handler: self.withItem(self.rebaseInteractiveOntoUpstream),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.upstreamStoredLocally),
			),
			Description:     "Interactive rebase onto upstream",
			DescriptionFunc: self.upstreamRebaseDescriptionFunc("Interactive rebase onto upstream"),
			Tooltip:         self.c.Tr.InteractiveRebaseTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseUpstreamOntoBase),
			Handler:           self.rebaseOntoBaseBranch,
			GetDisabledReason: self.require(self.singleItemSelected(self.canRebaseOntoBase)),
			Description:       "Rebase onto base branch",
			DescriptionFunc:   self.rebaseOntoBaseDescriptionFunc(),
			Tooltip:           self.c.Tr.RebaseOntoBaseBranchTooltip,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler: self.withItem(func(selectedBranch *models.Branch) error {
				return self.c.Helpers().Diff.OpenDiffToolForRef(selectedBranch)
			}),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
	}
}

func (self *BranchesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			branch := self.context().GetSelected()
			if branch == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoBranchesThisRepo)
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(branch.FullRefName())

				ptyTask := types.NewRunPtyTask(cmdObj.GetCmd())
				task = ptyTask

				pr, ok := self.c.Model().PullRequestsMap[branch.Name]
				if ok && presentation.ShouldShowPrForBranch(pr, branch.Name, self.c.UserConfig()) {
					icon := lo.Ternary(icons.IsIconEnabled(), icons.IconForRemoteUrl(pr.Url)+"  ", "")
					ptyTask.Prefix = style.PrintHyperlink(fmt.Sprintf("%s%s  %s  %s\n",
						icon,
						coloredStateText(pr.State),
						pr.Title,
						style.FgCyan.Sprintf("#%d", pr.Number)),
						pr.Url)
					ptyTask.Prefix += strings.Repeat("─", self.c.Contexts().Normal.GetView().InnerWidth()) + "\n"
				}
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: self.c.Tr.LogTitle,
					Task:  task,
				},
			})
		})
	}
}

func stateText(state string) string {
	var icon, label string
	switch state {
	case "OPEN":
		icon, label = " ", "Open"
	case "CLOSED":
		icon, label = " ", "Closed"
	case "MERGED":
		icon, label = " ", "Merged"
	case "DRAFT":
		icon, label = " ", "Draft"
	default:
		return ""
	}
	if icons.IsIconEnabled() {
		return icon + label
	}
	return label
}

func coloredStateText(state string) string {
	if icons.IsIconEnabled() {
		return fmt.Sprintf("%s%s%s",
			presentation.WithPrColor(state, "", false),
			presentation.WithPrColor(state, color.RGB(0xFF, 0xFF, 0xFF, false).Sprint(stateText(state)), true),
			presentation.WithPrColor(state, "", false))
	}

	return presentation.WithPrColor(state, stateText(state), false)
}

func (self *BranchesController) viewDivergenceFromUpstream(branch *models.Branch) error {
	upstream := lo.Ternary(branch.RemoteBranchStoredLocally(),
		branch.ShortUpstreamRefName(),
		self.c.Tr.UpstreamGenericName)

	return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
		Ref:                     branch,
		TitleRef:                fmt.Sprintf("%s <-> %s", branch.RefName(), upstream),
		RefToShowDivergenceFrom: branch.FullUpstreamRefName(),
		Context:                 self.context(),
		ShowBranchHeads:         false,
	})
}

func (self *BranchesController) viewDivergenceFromBaseBranch(branch *models.Branch) error {
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(branch, self.c.Model().MainBranches)
	if err != nil {
		return err
	}
	if baseBranch == "" {
		return errors.New(self.c.Tr.CouldNotDetermineBaseBranch)
	}
	shortBaseBranchName := helpers.ShortBranchName(baseBranch)

	return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
		Ref:                     branch,
		TitleRef:                fmt.Sprintf("%s <-> %s", branch.RefName(), shortBaseBranchName),
		RefToShowDivergenceFrom: baseBranch,
		Context:                 self.context(),
		ShowBranchHeads:         false,
	})
}

func (self *BranchesController) setUpstream(branch *models.Branch) error {
	return self.c.Helpers().Upstream.PromptForUpstreamWithoutInitialContent(branch, func(upstream string) error {
		upstreamRemote, upstreamBranch, err := self.c.Helpers().Upstream.ParseUpstream(upstream)
		if err != nil {
			return err
		}

		if err := self.c.Git().Branch.SetUpstream(upstreamRemote, upstreamBranch, branch.Name); err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC,
			Scope: []types.RefreshableView{
				types.BRANCHES,
				types.COMMITS,
			},
		})
		return nil
	})
}

func (self *BranchesController) unsetUpstream(branch *models.Branch) error {
	if err := self.c.Git().Branch.UnsetUpstream(branch.Name); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC,
		Scope: []types.RefreshableView{
			types.BRANCHES,
			types.COMMITS,
		},
	})
	return nil
}

func (self *BranchesController) upstreamStoredLocally(branch *models.Branch) *types.DisabledReason {
	if !branch.RemoteBranchStoredLocally() {
		return &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
	}
	return nil
}

func (self *BranchesController) branchIsTrackingRemote(branch *models.Branch) *types.DisabledReason {
	if !branch.IsTrackingRemote() {
		return &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
	}
	return nil
}

func (self *BranchesController) canViewDivergenceFromBase(branch *models.Branch) *types.DisabledReason {
	if self.c.Git() == nil {
		return nil
	}
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(branch, self.c.Model().MainBranches)
	if err != nil || baseBranch == "" {
		return &types.DisabledReason{Text: self.c.Tr.CouldNotDetermineBaseBranch}
	}
	return nil
}

// Returns "" on the cheatsheet path (no git client) or when no branch is
// selected.
func (self *BranchesController) upstreamShortName() string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil || !branch.RemoteBranchStoredLocally() {
		return ""
	}
	return branch.ShortUpstreamRefName()
}

func (self *BranchesController) baseBranchShortName() string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil {
		return ""
	}
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(branch, self.c.Model().MainBranches)
	if err != nil || baseBranch == "" {
		return ""
	}
	return helpers.ShortBranchName(baseBranch)
}

func (self *BranchesController) viewDivergenceFromUpstreamDescriptionFunc() func() string {
	return func() string {
		upstream := self.upstreamShortName()
		if upstream == "" {
			return self.c.Tr.ViewDivergenceFromUpstream
		}
		return fmt.Sprintf("%s '%s'", self.c.Tr.ViewDivergenceFromUpstream, upstream)
	}
}

func (self *BranchesController) viewDivergenceFromBaseDescriptionFunc() func() string {
	return func() string {
		base := self.baseBranchShortName()
		if base == "" {
			return "View divergence from base branch"
		}
		return fmt.Sprintf("%s '%s'", "View divergence from base branch", base)
	}
}

func (self *BranchesController) upstreamResetDescriptionFunc(staticLabel string) func() string {
	return func() string {
		upstream := self.upstreamShortName()
		if upstream == "" {
			return staticLabel
		}
		return fmt.Sprintf("%s '%s'", staticLabel, upstream)
	}
}

func (self *BranchesController) upstreamRebaseDescriptionFunc(staticLabel string) func() string {
	return func() string {
		upstream := self.upstreamShortName()
		if upstream == "" {
			return staticLabel
		}
		return fmt.Sprintf("%s '%s'", staticLabel, upstream)
	}
}

func (self *BranchesController) rebaseOntoBaseDescriptionFunc() func() string {
	return func() string {
		base := self.baseBranchShortName()
		if base == "" {
			return "Rebase onto base branch"
		}
		return utils.ResolvePlaceholderString(
			self.c.Tr.RebaseOntoBaseBranch,
			map[string]string{"baseBranch": base},
		)
	}
}

func (self *BranchesController) selectedBranchShortName() string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil {
		return ""
	}
	return branch.Name
}

func (self *BranchesController) rebaseBranchSimpleDescriptionFunc() func() string {
	return func() string {
		ref := self.selectedBranchShortName()
		if ref == "" {
			return "Simple rebase"
		}
		return utils.ResolvePlaceholderString(
			self.c.Tr.SimpleRebase,
			map[string]string{"ref": ref},
		)
	}
}

func (self *BranchesController) rebaseBranchInteractiveDescriptionFunc() func() string {
	return func() string {
		ref := self.selectedBranchShortName()
		if ref == "" {
			return "Interactive rebase"
		}
		return utils.ResolvePlaceholderString(
			self.c.Tr.InteractiveRebase,
			map[string]string{"ref": ref},
		)
	}
}

func (self *BranchesController) nonFastForwardMergeApplicable(branch *models.Branch) *types.DisabledReason {
	return self.c.Helpers().MergeAndRebase.NonFastForwardMergeDisabledReason(branch.Name)
}

func (self *BranchesController) fastForwardOnlyMergeApplicable(branch *models.Branch) *types.DisabledReason {
	return self.c.Helpers().MergeAndRebase.FastForwardOnlyMergeDisabledReason(branch.Name)
}

func (self *BranchesController) upstreamResetPreview(strength string) string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil || !branch.RemoteBranchStoredLocally() {
		return ""
	}
	return style.FgRed.Sprintf("reset --%s %s", strength, branch.ShortUpstreamRefName())
}

func (self *BranchesController) gitMixedResetToUpstream(branch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(branch.ShortUpstreamRefName(), branch.FullUpstreamRefName(), "mixed")
}

func (self *BranchesController) gitSoftResetToUpstream(branch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(branch.ShortUpstreamRefName(), branch.FullUpstreamRefName(), "soft")
}

func (self *BranchesController) gitHardResetToUpstream(branch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(branch.ShortUpstreamRefName(), branch.FullUpstreamRefName(), "hard")
}

func (self *BranchesController) rebaseSimpleOntoUpstream(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(branch.ShortUpstreamRefName(), helpers.RebaseVariantSimple)
}

func (self *BranchesController) rebaseInteractiveOntoUpstream(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(branch.ShortUpstreamRefName(), helpers.RebaseVariantInteractive)
}

func (self *BranchesController) Context() types.Context {
	return self.context()
}

func (self *BranchesController) context() *context.BranchesContext {
	return self.c.Contexts().Branches
}

func (self *BranchesController) press(selectedBranch *models.Branch) error {
	if selectedBranch == self.c.Helpers().Refs.GetCheckedOutRef() {
		return errors.New(self.c.Tr.AlreadyCheckedOutBranch)
	}

	worktreeForRef, ok := self.worktreeForBranch(selectedBranch)
	if ok && !worktreeForRef.IsCurrent {
		return self.promptToCheckoutWorktree(worktreeForRef)
	}

	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	return self.c.Helpers().Refs.CheckoutRef(selectedBranch.Name, types.CheckoutRefOptions{})
}

func (self *BranchesController) notPulling() *types.DisabledReason {
	currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	if currentBranch != nil {
		op := self.c.State().GetItemOperation(currentBranch)
		if op == types.ItemOperationFastForwarding || op == types.ItemOperationPulling {
			return &types.DisabledReason{Text: self.c.Tr.CantCheckoutBranchWhilePulling}
		}
	}

	return nil
}

func (self *BranchesController) worktreeForBranch(branch *models.Branch) (*models.Worktree, bool) {
	return git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees)
}

func (self *BranchesController) promptToCheckoutWorktree(worktree *models.Worktree) error {
	prompt := utils.ResolvePlaceholderString(self.c.Tr.AlreadyCheckedOutByWorktree, map[string]string{
		"worktreeName": worktree.Name,
	})

	return self.c.ConfirmIf(!self.c.UserConfig().Gui.SkipSwitchWorktreeOnCheckoutWarning, types.ConfirmOpts{
		Title:  self.c.Tr.SwitchToWorktree,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.Helpers().Worktree.Switch(worktree, context.LOCAL_BRANCHES_CONTEXT_KEY)
		},
	})
}

func (self *BranchesController) handleCreatePullRequest(selectedBranch *models.Branch) error {
	if !selectedBranch.IsTrackingRemote() {
		return errors.New(self.c.Tr.PullRequestNoUpstream)
	}
	return self.createPullRequest(selectedBranch.UpstreamBranch, "")
}

func (self *BranchesController) handleCreatePullRequestMenu(selectedBranch *models.Branch) error {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()

	return self.createPullRequestMenu(selectedBranch, checkedOutBranch)
}

func (self *BranchesController) getPullRequestURL() (string, error) {
	branch := self.context().GetSelected()
	if pr, ok := self.c.Model().PullRequestsMap[branch.Name]; ok {
		return pr.Url, nil
	}

	branchExistsOnRemote := self.c.Git().Remote.CheckRemoteBranchExists(branch.Name)

	if !branchExistsOnRemote {
		return "", errors.New(self.c.Tr.NoBranchOnRemote)
	}

	return self.c.Helpers().Host.GetPullRequestURL(branch.Name, "")
}

func (self *BranchesController) copyPullRequestURL() error {
	url, err := self.getPullRequestURL()
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.CopyPullRequestURL)
	if err := self.c.OS().CopyToClipboard(url); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (self *BranchesController) forceCheckout() error {
	branch := self.context().GetSelected()
	message := self.c.Tr.SureForceCheckout
	title := self.c.Tr.ForceCheckoutBranch

	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ForceCheckoutBranch)
			if err := self.c.Git().Branch.Checkout(branch.Name, git_commands.CheckoutOptions{Force: true}); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return nil
		},
	})

	return nil
}

func (self *BranchesController) checkoutPreviousBranch() error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	return self.c.Helpers().Refs.CheckoutPreviousRef()
}

func (self *BranchesController) checkoutByName() error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.BranchName + ":",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRefsSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			self.c.LogAction("Checkout branch")
			_, branchName, found := self.c.Helpers().Refs.ParseRemoteBranchName(response)
			if found {
				return self.c.Helpers().Refs.CheckoutRemoteBranch(response, branchName)
			}
			return self.c.Helpers().Refs.CheckoutRef(response, types.CheckoutRefOptions{
				OnRefNotFound: func(ref string) error {
					self.c.Confirm(types.ConfirmOpts{
						Title:  self.c.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", self.c.Tr.BranchNotFoundPrompt, ref, "?"),
						HandleConfirm: func() error {
							return self.createNewBranchWithName(ref)
						},
					})

					return nil
				},
			})
		},
	},
	)

	return nil
}

func (self *BranchesController) createNewBranchWithName(newBranchName string) error {
	branch := self.context().GetSelected()
	if branch == nil {
		return nil
	}

	if err := self.c.Git().Branch.New(newBranchName, branch.FullRefName()); err != nil {
		return err
	}

	self.c.Helpers().Refs.SelectFirstBranchAndFirstCommit()
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, KeepBranchSelectionIndex: true})
	return nil
}

func (self *BranchesController) localDelete(branches []*models.Branch) error {
	return self.c.Helpers().BranchesHelper.ConfirmLocalDelete(branches)
}

func (self *BranchesController) remoteDelete(branches []*models.Branch) error {
	remoteBranches := lo.Map(branches, func(branch *models.Branch, _ int) *models.RemoteBranch {
		return &models.RemoteBranch{Name: branch.UpstreamBranch, RemoteName: branch.UpstreamRemote}
	})
	return self.c.Helpers().BranchesHelper.ConfirmDeleteRemote(remoteBranches, false)
}

func (self *BranchesController) localAndRemoteDelete(branches []*models.Branch) error {
	return self.c.Helpers().BranchesHelper.ConfirmLocalAndRemoteDelete(branches)
}

func (self *BranchesController) notDeletingCheckedOutBranch(branches []*models.Branch, _ int, _ int) *types.DisabledReason {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	if checkedOutBranch == nil {
		return nil
	}
	isBranchCheckedOut := lo.SomeBy(branches, func(branch *models.Branch) bool {
		return checkedOutBranch.Name == branch.Name
	})
	if isBranchCheckedOut {
		return &types.DisabledReason{Text: self.c.Tr.CantDeleteCheckOutBranch}
	}
	return nil
}

func (self *BranchesController) allBranchesHaveUpstream(branches []*models.Branch, _ int, _ int) *types.DisabledReason {
	hasUpstream := lo.EveryBy(branches, func(branch *models.Branch) bool {
		return branch.IsTrackingRemote() && !branch.UpstreamGone
	})
	if !hasUpstream {
		return &types.DisabledReason{
			Text: lo.Ternary(len(branches) > 1, self.c.Tr.UpstreamsNotSetError, self.c.Tr.UpstreamNotSetError),
		}
	}
	return nil
}

func (self *BranchesController) deleteBranchDescriptionFunc(singular string, plural string) func() string {
	return func() string {
		items, _, _ := self.context().GetSelectedItems()
		return lo.Ternary(len(items) > 1, plural, singular)
	}
}

func (self *BranchesController) mergeRegular(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(branch.Name, git_commands.MERGE_VARIANT_REGULAR)
}

func (self *BranchesController) mergeNonFastForward(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(branch.Name, git_commands.MERGE_VARIANT_NON_FAST_FORWARD)
}

func (self *BranchesController) mergeFastForward(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformMerge(branch.Name, git_commands.MERGE_VARIANT_FAST_FORWARD)
}

func (self *BranchesController) mergeSquash(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformSquashMerge(branch.Name)
}

func (self *BranchesController) mergeSquashCommitted(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformSquashMergeCommitted(branch.Name)
}

func (self *BranchesController) rebaseSimple(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(branch.Name, helpers.RebaseVariantSimple)
}

func (self *BranchesController) rebaseInteractive(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef(branch.Name, helpers.RebaseVariantInteractive)
}

func (self *BranchesController) rebaseOntoBaseBranch() error {
	return self.c.Helpers().MergeAndRebase.PerformRebaseOntoRef("", helpers.RebaseVariantOntoBase)
}

func (self *BranchesController) notRebasingOntoSelf(branch *models.Branch) *types.DisabledReason {
	if len(self.c.Model().Branches) == 0 {
		return nil
	}
	if branch.Name == self.c.Model().Branches[0].Name {
		return &types.DisabledReason{Text: self.c.Tr.CantRebaseOntoSelf}
	}
	return nil
}

func (self *BranchesController) canRebaseOntoBase(_ *models.Branch) *types.DisabledReason {
	if self.c.Git() == nil {
		return nil
	}
	_, reason := self.c.Helpers().MergeAndRebase.RebaseOntoBaseBranchName()
	return reason
}

func (self *BranchesController) fastForward(branch *models.Branch) error {
	if !branch.IsTrackingRemote() {
		return errors.New(self.c.Tr.FwdNoUpstream)
	}
	if !branch.RemoteBranchStoredLocally() {
		return errors.New(self.c.Tr.FwdNoLocalUpstream)
	}
	if branch.IsAheadForPull() {
		return errors.New(self.c.Tr.FwdCommitsToPush)
	}

	action := self.c.Tr.Actions.FastForwardBranch

	return self.c.WithInlineStatus(branch, types.ItemOperationFastForwarding, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		worktree, ok := self.worktreeForBranch(branch)
		if ok {
			self.c.LogAction(action)

			worktreeGitDir := ""
			worktreePath := ""
			// if it is the current worktree path, no need to specify the path
			if !worktree.IsCurrent {
				worktreeGitDir = worktree.GitDir
				worktreePath = worktree.Path
			}

			err := self.c.Git().Sync.Pull(
				task,
				git_commands.PullOptions{
					RemoteName:      branch.UpstreamRemote,
					BranchName:      branch.UpstreamBranch,
					FastForwardOnly: true,
					WorktreeGitDir:  worktreeGitDir,
					WorktreePath:    worktreePath,
				},
			)
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return err
		}

		self.c.LogAction(action)

		err := self.c.Git().Sync.FastForward(
			task, branch.Name, branch.UpstreamRemote, branch.UpstreamBranch,
		)
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		return err
	})
}

func (self *BranchesController) createTag(branch *models.Branch) error {
	return self.c.Helpers().Tags.OpenCreateTagPrompt(branch.FullRefName(), func() {})
}

func (self *BranchesController) createSortMenu() error {
	return self.c.Helpers().Refs.CreateSortOrderMenu(
		[]string{"recency", "alphabetical", "date"},
		self.c.Tr.SortOrderPromptLocalBranches,
		func(sortOrder string) error {
			if self.c.UserConfig().Git.LocalBranchSortOrder != sortOrder {
				self.c.UserConfig().Git.LocalBranchSortOrder = sortOrder
				self.c.Contexts().Branches.SetSelection(0)
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
				return nil
			}
			return nil
		},
		self.c.UserConfig().Git.LocalBranchSortOrder)
}

func (self *BranchesController) gitMixedResetToRef(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.Name, selectedBranch.FullRefName(), "mixed")
}

func (self *BranchesController) gitSoftResetToRef(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.Name, selectedBranch.FullRefName(), "soft")
}

func (self *BranchesController) gitHardResetToRef(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.PerformGitReset(selectedBranch.Name, selectedBranch.FullRefName(), "hard")
}

func (self *BranchesController) gitResetPreview(s style.TextStyle, strength string) string {
	if self.c.Git() == nil {
		return ""
	}
	branch := self.context().GetSelected()
	if branch == nil {
		return ""
	}
	return s.Sprintf("reset --%s %s", strength, branch.Name)
}

func (self *BranchesController) rename(branch *models.Branch) error {
	promptForNewName := func() error {
		self.c.Prompt(types.PromptOpts{
			Title:          self.c.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				self.c.LogAction(self.c.Tr.Actions.RenameBranch)
				if err := self.c.Git().Branch.Rename(branch.Name, helpers.SanitizedBranchName(newBranchName)); err != nil {
					return err
				}

				// need to find where the branch is now so that we can re-select it. That means we need to refetch the branches synchronously and then find our branch
				self.c.Refresh(types.RefreshOptions{
					Mode:  types.SYNC,
					Scope: []types.RefreshableView{types.BRANCHES, types.WORKTREES},
				})

				// now that we've got our stuff again we need to find that branch and reselect it.
				for i, newBranch := range self.c.Model().Branches {
					if newBranch.Name == newBranchName {
						self.context().SetSelection(i)
						self.context().HandleRender()
					}
				}

				return nil
			},
		})

		return nil
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	return self.c.ConfirmIf(branch.IsTrackingRemote(), types.ConfirmOpts{
		Title:         self.c.Tr.RenameBranch,
		Prompt:        self.c.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})
}

func (self *BranchesController) newBranch(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.NewBranch(selectedBranch.FullRefName(), selectedBranch.RefName(), "")
}

func (self *BranchesController) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*types.MenuItem, 0, 4)

	fromToLabelColumns := func(from string, to string) []string {
		return []string{fmt.Sprintf("%s → %s", from, to)}
	}

	menuItemsForBranch := func(branch *models.Branch) []*types.MenuItem {
		return []*types.MenuItem{
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.DefaultBranch),
				OnPress: func() error {
					return self.handleCreatePullRequest(branch)
				},
			},
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.SelectBranch),
				OnPress: func() error {
					if !branch.IsTrackingRemote() {
						return errors.New(self.c.Tr.PullRequestNoUpstream)
					}

					if len(self.c.Model().Remotes) == 1 {
						toRemote := self.c.Model().Remotes[0].Name
						self.c.Log.Debugf("PR will target the only existing remote '%s'", toRemote)
						return self.promptForTargetBranchNameAndCreatePullRequest(branch, toRemote)
					}

					self.c.Prompt(types.PromptOpts{
						Title:               self.c.Tr.SelectTargetRemote,
						FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
						HandleConfirm: func(toRemote string) error {
							self.c.Log.Debugf("PR will target remote '%s'", toRemote)

							return self.promptForTargetBranchNameAndCreatePullRequest(branch, toRemote)
						},
					})

					return nil
				},
			},
		}
	}

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems,
			&types.MenuItem{
				LabelColumns: fromToLabelColumns(checkedOutBranch.Name, selectedBranch.Name),
				OnPress: func() error {
					if !checkedOutBranch.IsTrackingRemote() || !selectedBranch.IsTrackingRemote() {
						return errors.New(self.c.Tr.PullRequestNoUpstream)
					}
					return self.createPullRequest(checkedOutBranch.UpstreamBranch, selectedBranch.UpstreamBranch)
				},
			},
		)
		menuItems = append(menuItems, menuItemsForBranch(checkedOutBranch)...)
	}

	menuItems = append(menuItems, menuItemsForBranch(selectedBranch)...)

	return self.c.Menu(types.CreateMenuOptions{Title: fmt.Sprint(self.c.Tr.CreatePullRequestOptions), Items: menuItems})
}

func (self *BranchesController) promptForTargetBranchNameAndCreatePullRequest(fromBranch *models.Branch, toRemote string) error {
	remoteDoesNotExist := lo.NoneBy(self.c.Model().Remotes, func(remote *models.Remote) bool {
		return remote.Name == toRemote
	})
	if remoteDoesNotExist {
		return fmt.Errorf(self.c.Tr.NoValidRemoteName, toRemote)
	}

	self.c.Prompt(types.PromptOpts{
		Title:               fmt.Sprintf("%s → %s/", fromBranch.UpstreamBranch, toRemote),
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteBranchesForRemoteSuggestionsFunc(toRemote),
		HandleConfirm: func(toBranch string) error {
			self.c.Log.Debugf("PR will target branch '%s' on remote '%s'", toBranch, toRemote)
			return self.createPullRequest(fromBranch.UpstreamBranch, toBranch)
		},
	})

	return nil
}

func (self *BranchesController) createPullRequest(from string, to string) error {
	url, err := self.c.Helpers().Host.GetPullRequestURL(from, to)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)

	if err := self.c.OS().OpenLink(url); err != nil {
		return err
	}

	return nil
}

func (self *BranchesController) branchIsReal(branch *models.Branch) *types.DisabledReason {
	if !branch.IsRealBranch() {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemIsNotABranch}
	}

	return nil
}

func (self *BranchesController) branchHasPR(branch *models.Branch) *types.DisabledReason {
	if _, ok := self.c.Model().PullRequestsMap[branch.Name]; !ok {
		return &types.DisabledReason{Text: self.c.Tr.NoPullRequestForBranch, ShowErrorInPanel: true}
	}

	return nil
}

func (self *BranchesController) openPRInBrowser(branch *models.Branch) error {
	pr, ok := self.c.Model().PullRequestsMap[branch.Name]
	if !ok {
		// Should be guarded against by the DisabledReason check, but be defensive in case
		// PullRequestsMap was updated concurrently by a background refresh
		return errors.New(self.c.Tr.NoPullRequestForBranch)
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)

	return self.c.OS().OpenLink(pr.Url)
}

func (self *BranchesController) branchesAreReal(selectedBranches []*models.Branch, startIdx int, endIdx int) *types.DisabledReason {
	if !lo.EveryBy(selectedBranches, func(branch *models.Branch) bool {
		return branch.IsRealBranch()
	}) {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemIsNotABranch}
	}

	return nil
}

func (self *BranchesController) notMergingIntoYourself(branch *models.Branch) *types.DisabledReason {
	selectedBranchName := branch.Name
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef().Name

	if checkedOutBranch == selectedBranchName {
		return &types.DisabledReason{Text: self.c.Tr.CantMergeBranchIntoItself}
	}

	return nil
}
