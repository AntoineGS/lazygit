package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type GitFlowController struct {
	baseController
	*ListControllerTrait[*models.Branch]
	c *ControllerCommon
}

var _ types.IController = &GitFlowController{}

func NewGitFlowController(
	c *ControllerCommon,
) *GitFlowController {
	return &GitFlowController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().Branches,
			c.Contexts().Branches.GetSelected,
			c.Contexts().Branches.GetSelectedItems,
		),
		c: c,
	}
}

func (self *GitFlowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	gitFlowDisabledReason := func() *types.DisabledReason {
		if self.c.Git() == nil {
			return nil
		}
		if !self.c.Git().Flow.GitFlowEnabled() {
			return &types.DisabledReason{Text: "You need to install git-flow and enable it in this repo to use git-flow features"}
		}
		return nil
	}

	return []*types.Binding{
		{
			Key:             opts.GetKey(opts.Config.Branches.GitFlowFinish),
			Handler:         self.withItem(self.gitFlowFinish),
			Description:     "Finish git-flow branch",
			DescriptionFunc: self.gitFlowFinishLabel,
			GetDisabledReason: func() *types.DisabledReason {
				if r := gitFlowDisabledReason(); r != nil {
					return r
				}
				return self.require(self.singleItemSelected())()
			},
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.GitFlowStartFeature),
			Handler:           self.startGitFlowFeature,
			Description:       "Start git-flow feature",
			GetDisabledReason: gitFlowDisabledReason,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.GitFlowStartHotfix),
			Handler:           self.startGitFlowHotfix,
			Description:       "Start git-flow hotfix",
			GetDisabledReason: gitFlowDisabledReason,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.GitFlowStartBugfix),
			Handler:           self.startGitFlowBugfix,
			Description:       "Start git-flow bugfix",
			GetDisabledReason: gitFlowDisabledReason,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.GitFlowStartRelease),
			Handler:           self.startGitFlowRelease,
			Description:       "Start git-flow release",
			GetDisabledReason: gitFlowDisabledReason,
		},
	}
}

func (self *GitFlowController) gitFlowFinishLabel() string {
	if self.c.Git() == nil {
		return "Finish git-flow branch"
	}
	branch := self.c.Contexts().Branches.GetSelected()
	if branch == nil {
		return "Finish git-flow branch"
	}
	return fmt.Sprintf("finish branch '%s'", branch.Name)
}

func (self *GitFlowController) gitFlowFinish(branch *models.Branch) error {
	cmdObj, err := self.c.Git().Flow.FinishCmdObj(branch.Name)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.GitFlowFinish)
	return self.c.RunSubprocessAndRefresh(cmdObj)
}

func (self *GitFlowController) startGitFlowFeature() error {
	return self.startGitFlowBranch("feature")
}

func (self *GitFlowController) startGitFlowHotfix() error {
	return self.startGitFlowBranch("hotfix")
}

func (self *GitFlowController) startGitFlowBugfix() error {
	return self.startGitFlowBranch("bugfix")
}

func (self *GitFlowController) startGitFlowRelease() error {
	return self.startGitFlowBranch("release")
}

func (self *GitFlowController) startGitFlowBranch(branchType string) error {
	title := utils.ResolvePlaceholderString(self.c.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

	self.c.Prompt(types.PromptOpts{
		Title: title,
		HandleConfirm: func(name string) error {
			self.c.LogAction(self.c.Tr.Actions.GitFlowStart)
			return self.c.RunSubprocessAndRefresh(
				self.c.Git().Flow.StartCmdObj(branchType, name),
			)
		},
	})

	return nil
}
