package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type MergeAndRebaseHelper struct {
	c *HelperCommon
}

func NewMergeAndRebaseHelper(
	c *HelperCommon,
) *MergeAndRebaseHelper {
	return &MergeAndRebaseHelper{
		c: c,
	}
}

type RebaseOption string

const (
	REBASE_OPTION_CONTINUE string = "continue"
	REBASE_OPTION_ABORT    string = "abort"
	REBASE_OPTION_SKIP     string = "skip"
)

func (self *MergeAndRebaseHelper) ContinueRebase() error {
	return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
}

func (self *MergeAndRebaseHelper) AbortRebase() error {
	return self.genericMergeCommand(REBASE_OPTION_ABORT)
}

func (self *MergeAndRebaseHelper) SkipRebase() error {
	return self.genericMergeCommand(REBASE_OPTION_SKIP)
}

func (self *MergeAndRebaseHelper) genericMergeCommand(command string) error {
	status := self.c.Git().Status.WorkingTreeState()

	if status.None() {
		return errors.New(self.c.Tr.NotMergingOrRebasing)
	}

	self.c.LogAction(fmt.Sprintf("Merge/Rebase: %s", command))
	effectiveStatus := status.Effective()
	if effectiveStatus == models.WORKING_TREE_STATE_REBASING {
		todoFile, err := os.ReadFile(
			filepath.Join(self.c.Git().RepoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		)

		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			self.c.LogCommand(string(todoFile), false)
		}
	}

	commandType := status.CommandName()

	// we should end up with a command like 'git merge --continue'

	// it's impossible for a rebase to require a commit so we'll use a subprocess only if it's a merge
	needsSubprocess := (effectiveStatus == models.WORKING_TREE_STATE_MERGING && command != REBASE_OPTION_ABORT && self.c.UserConfig().Git.Merging.ManualCommit) ||
		// but we'll also use a subprocess if we have exec todos; those are likely to be lengthy build
		// tasks whose output the user will want to see in the terminal
		(effectiveStatus == models.WORKING_TREE_STATE_REBASING && command != REBASE_OPTION_ABORT && self.hasExecTodos())

	if needsSubprocess {
		// TODO: see if we should be calling more of the code from self.Git.Rebase.GenericMergeOrRebaseAction
		return self.c.RunSubprocessAndRefresh(
			self.c.Git().Rebase.GenericMergeOrRebaseActionCmdObj(commandType, command),
		)
	}
	result := self.c.Git().Rebase.GenericMergeOrRebaseAction(commandType, command)
	if err := self.CheckMergeOrRebase(result); err != nil {
		return err
	}
	return nil
}

func (self *MergeAndRebaseHelper) hasExecTodos() bool {
	for _, commit := range self.c.Model().Commits {
		if !commit.IsTODO() {
			break
		}
		if commit.Action == todo.Exec {
			return true
		}
	}
	return false
}

var conflictStrings = []string{
	"Failed to merge in the changes",
	"When you have resolved this problem",
	"fix conflicts",
	"Resolve all conflicts manually",
	"Merge conflict in file",
	"hint: after resolving the conflicts",
	"CONFLICT (content):",
}

func isMergeConflictErr(errStr string) bool {
	for _, str := range conflictStrings {
		if strings.Contains(errStr, str) {
			return true
		}
	}

	return false
}

func (self *MergeAndRebaseHelper) CheckMergeOrRebaseWithRefreshOptions(result error, refreshOptions types.RefreshOptions) error {
	self.c.Refresh(refreshOptions)

	if result == nil {
		return nil
	} else if strings.Contains(result.Error(), "No changes - did you forget to use") {
		return self.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "The previous cherry-pick is now empty") {
		return self.genericMergeCommand(REBASE_OPTION_SKIP)
	} else if strings.Contains(result.Error(), "No rebase in progress?") {
		// assume in this case that we're already done
		return nil
	}
	return self.CheckForConflicts(result)
}

func (self *MergeAndRebaseHelper) CheckMergeOrRebase(result error) error {
	return self.CheckMergeOrRebaseWithRefreshOptions(result, types.RefreshOptions{Mode: types.ASYNC})
}

func (self *MergeAndRebaseHelper) CheckForConflicts(result error) error {
	if result == nil {
		return nil
	}

	if isMergeConflictErr(result.Error()) {
		return self.PromptForConflictHandling()
	}

	return result
}

func (self *MergeAndRebaseHelper) PromptForConflictHandling() error {
	mode := self.c.Git().Status.WorkingTreeState().CommandName()
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.FoundConflictsTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.ViewConflictsMenuItem,
				OnPress: func() error {
					self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
					return nil
				},
			},
			{
				Label: fmt.Sprintf(self.c.Tr.AbortMenuItem, mode),
				OnPress: func() error {
					return self.genericMergeCommand(REBASE_OPTION_ABORT)
				},
				Key: gocui.NewKeyRune('a'),
			},
		},
		HideCancel: true,
	})
}

func (self *MergeAndRebaseHelper) AbortMergeOrRebaseWithConfirm() error {
	// prompt user to confirm that they want to abort, then do it
	mode := self.c.Git().Status.WorkingTreeState().CommandName()
	self.c.Confirm(types.ConfirmOpts{
		Title:  fmt.Sprintf(self.c.Tr.AbortTitle, mode),
		Prompt: fmt.Sprintf(self.c.Tr.AbortPrompt, mode),
		HandleConfirm: func() error {
			return self.genericMergeCommand(REBASE_OPTION_ABORT)
		},
	})

	return nil
}

// PromptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (self *MergeAndRebaseHelper) PromptToContinueRebase() error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Continue,
		Prompt: fmt.Sprintf(self.c.Tr.ConflictsResolved, self.c.Git().Status.WorkingTreeState().CommandName()),
		HandleConfirm: func() error {
			// By the time we get here, we might have unstaged changes again,
			// e.g. if the user had to fix build errors after resolving the
			// conflicts, but after lazygit opened the prompt already. Ask again
			// to auto-stage these.

			// Need to refresh the files to be really sure if this is the case.
			// We would otherwise be relying on lazygit's auto-refresh on focus,
			// but this is not supported by all terminals or on all platforms.
			self.c.Refresh(types.RefreshOptions{
				Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES},
			})

			unstagedFiles := GetUnstagedFilesExceptSubmodules(self.c.Model().Files, self.c.Model().Submodules)
			if len(unstagedFiles) > 0 {
				self.c.Confirm(types.ConfirmOpts{
					Title:  self.c.Tr.Continue,
					Prompt: self.c.Tr.UnstagedFilesAfterConflictsResolved,
					HandleConfirm: func() error {
						self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
						if err := self.c.Git().WorkingTree.StageFiles(unstagedFiles, []string{}); err != nil {
							return err
						}

						return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
					},
				})

				return nil
			}

			return self.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})

	return nil
}

type RebaseVariant string

const (
	RebaseVariantSimple      RebaseVariant = "simple"
	RebaseVariantInteractive RebaseVariant = "interactive"
	RebaseVariantOntoBase    RebaseVariant = "onto-base"
)

// For RebaseVariantOntoBase, ref is ignored and the base branch is
// resolved internally.
func (self *MergeAndRebaseHelper) PerformRebaseOntoRef(ref string, variant RebaseVariant) error {
	checkedOutBranch := self.c.Model().Branches[0]

	if variant == RebaseVariantOntoBase {
		baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(checkedOutBranch, self.c.Model().MainBranches)
		if err != nil {
			return err
		}
		if baseBranch == "" {
			return errors.New(self.c.Tr.CouldNotDetermineBaseBranch)
		}
		ref = baseBranch
	}

	if variant == RebaseVariantInteractive {
		self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
		baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
		var err error
		if baseCommit != "" {
			err = self.c.Git().Rebase.EditRebaseFromBaseCommit(ref, baseCommit)
		} else {
			err = self.c.Git().Rebase.EditRebase(ref)
		}
		if err = self.CheckMergeOrRebase(err); err != nil {
			return err
		}
		if err = self.ResetMarkedBaseCommit(); err != nil {
			return err
		}
		self.c.Context().Push(self.c.Contexts().LocalCommits, types.OnFocusOpts{})
		return nil
	}

	// Simple and OntoBase variants share the same body.
	self.c.LogAction(self.c.Tr.Actions.RebaseBranch)
	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(task gocui.Task) error {
		baseCommit := self.c.Modes().MarkedBaseCommit.GetHash()
		var err error
		if baseCommit != "" {
			err = self.c.Git().Rebase.RebaseBranchFromBaseCommit(ref, baseCommit)
		} else {
			err = self.c.Git().Rebase.RebaseBranch(ref)
		}
		err = self.CheckMergeOrRebase(err)
		if err == nil {
			return self.ResetMarkedBaseCommit()
		}
		return err
	})
}

func (self *MergeAndRebaseHelper) RebaseOntoBaseBranchName() (string, *types.DisabledReason) {
	if self.c.Git() == nil {
		return "", nil
	}
	checkedOutBranch := self.c.Model().Branches[0]
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(checkedOutBranch, self.c.Model().MainBranches)
	if err != nil || baseBranch == "" {
		return "", &types.DisabledReason{Text: self.c.Tr.CouldNotDetermineBaseBranch}
	}
	return baseBranch, nil
}

func (self *MergeAndRebaseHelper) PerformMerge(refName string, variant git_commands.MergeVariant) error {
	if self.c.Git().Branch.IsHeadDetached() {
		return errors.New("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := self.c.Model().Branches[0].Name
	if checkedOutBranchName == refName {
		return errors.New(self.c.Tr.CantMergeBranchIntoItself)
	}
	if variant == git_commands.MERGE_VARIANT_FAST_FORWARD && !self.c.Git().Branch.CanDoFastForwardMerge(refName) {
		// Surface the friendly error before letting git fail with its
		// own message.
		return errors.New(utils.ResolvePlaceholderString(
			self.c.Tr.CannotFastForwardMerge,
			map[string]string{
				"checkedOutBranch": checkedOutBranchName,
				"selectedBranch":   refName,
			},
		))
	}
	return self.RegularMerge(refName, variant)()
}

func (self *MergeAndRebaseHelper) PerformSquashMerge(refName string) error {
	if self.c.Git().Branch.IsHeadDetached() {
		return errors.New("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := self.c.Model().Branches[0].Name
	if checkedOutBranchName == refName {
		return errors.New(self.c.Tr.CantMergeBranchIntoItself)
	}
	return self.SquashMergeUncommitted(refName)()
}

// PerformSquashMergeCommitted runs a squash merge AND commits it.
func (self *MergeAndRebaseHelper) PerformSquashMergeCommitted(refName string) error {
	if self.c.Git().Branch.IsHeadDetached() {
		return errors.New("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := self.c.Model().Branches[0].Name
	if checkedOutBranchName == refName {
		return errors.New(self.c.Tr.CantMergeBranchIntoItself)
	}
	return self.SquashMergeCommitted(refName, checkedOutBranchName)()
}

func (self *MergeAndRebaseHelper) MergeIntoSelfDisabledReason(selectedRefName string) *types.DisabledReason {
	if self.c.Git() == nil {
		return nil
	}
	if self.c.Git().Branch.IsHeadDetached() {
		return &types.DisabledReason{Text: "Cannot merge while in detached HEAD state"}
	}
	if len(self.c.Model().Branches) == 0 {
		return nil
	}
	if self.c.Model().Branches[0].Name == selectedRefName {
		return &types.DisabledReason{Text: self.c.Tr.CantMergeBranchIntoItself}
	}
	return nil
}

func (self *MergeAndRebaseHelper) RegularMerge(refName string, variant git_commands.MergeVariant) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.Merge)
		err := self.c.Git().Branch.Merge(refName, variant)
		return self.CheckMergeOrRebase(err)
	}
}

func (self *MergeAndRebaseHelper) SquashMergeUncommitted(refName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		err := self.c.Git().Branch.Merge(refName, git_commands.MERGE_VARIANT_SQUASH)
		return self.CheckMergeOrRebase(err)
	}
}

func (self *MergeAndRebaseHelper) SquashMergeCommitted(refName, checkedOutBranchName string) func() error {
	return func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashMerge)
		err := self.c.Git().Branch.Merge(refName, git_commands.MERGE_VARIANT_SQUASH)
		if err = self.CheckMergeOrRebase(err); err != nil {
			return err
		}
		message := utils.ResolvePlaceholderString(self.c.UserConfig().Git.Merging.SquashMergeMessage, map[string]string{
			"selectedRef":   refName,
			"currentBranch": checkedOutBranchName,
		})
		err = self.c.Git().Commit.CommitCmdObj(message, "", false).Run()
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	}
}

// Returns wantsFastForward, wantsNonFastForward. These will never both be true, but they can both be false.
func (self *MergeAndRebaseHelper) fastForwardMergeUserPreference() (bool, bool) {
	// Check user config first, because it takes precedence over git config
	mergingArgs := self.c.UserConfig().Git.Merging.Args
	if strings.Contains(mergingArgs, "--ff") { // also covers "--ff-only"
		return true, false
	}

	if strings.Contains(mergingArgs, "--no-ff") {
		return false, true
	}

	// Then check git config
	mergeFfConfig := self.c.Git().Config.GetMergeFF()
	if mergeFfConfig == "true" || mergeFfConfig == "only" {
		return true, false
	}

	if mergeFfConfig == "false" {
		return false, true
	}

	return false, false
}

// Disabled when the regular merge already creates a merge commit — i.e.
// the user's config prevents fast-forward or the branches aren't
// fast-forwardable. (Adding --no-ff would be redundant.)
func (self *MergeAndRebaseHelper) NonFastForwardMergeDisabledReason(refName string) *types.DisabledReason {
	if self.c.Git() == nil {
		return nil
	}
	wantFF, wantNFF := self.fastForwardMergeUserPreference()
	canFF := self.c.Git().Branch.CanDoFastForwardMerge(refName)
	if wantNFF || (!wantFF && !canFF) {
		return &types.DisabledReason{Text: self.c.Tr.MergeNonFastForwardNotApplicable}
	}
	return nil
}

func (self *MergeAndRebaseHelper) FastForwardOnlyMergeDisabledReason(refName string) *types.DisabledReason {
	if self.c.Git() == nil {
		return nil
	}
	wantFF, wantNFF := self.fastForwardMergeUserPreference()
	canFF := self.c.Git().Branch.CanDoFastForwardMerge(refName)
	if !wantNFF && (wantFF || canFF) {
		return &types.DisabledReason{Text: self.c.Tr.MergeFastForwardNotApplicable}
	}
	if !canFF {
		checkedOutBranch := ""
		if len(self.c.Model().Branches) > 0 {
			checkedOutBranch = self.c.Model().Branches[0].Name
		}
		return &types.DisabledReason{
			Text: utils.ResolvePlaceholderString(
				self.c.Tr.CannotFastForwardMerge,
				map[string]string{
					"checkedOutBranch": checkedOutBranch,
					"selectedBranch":   refName,
				},
			),
			ShowErrorInPanel: true,
		}
	}
	return nil
}

func (self *MergeAndRebaseHelper) ResetMarkedBaseCommit() error {
	self.c.Modes().MarkedBaseCommit.Reset()
	self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
	return nil
}
