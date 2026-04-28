package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type AllBranchesLogHelper struct {
	c *HelperCommon
}

func NewAllBranchesLogHelper(c *HelperCommon) *AllBranchesLogHelper {
	return &AllBranchesLogHelper{c: c}
}

// ShowAllBranchLogs renders the configured all-branches log command into
// the main view. Body lifted verbatim from
// StatusController.showAllBranchLogs so non-status callers (e.g. universal
// bindings via GlobalController) can invoke it.
func (self *AllBranchesLogHelper) ShowAllBranchLogs() {
	cmdObj := self.c.Git().Branch.AllBranchesLogCmdObj()
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	title := self.c.Tr.LogTitle
	if i, n := self.c.Git().Branch.GetAllBranchesLogIdxAndCount(); n > 1 {
		title = fmt.Sprintf(self.c.Tr.LogXOfYTitle, i+1, n)
	}
	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: title,
			Task:  task,
		},
	})
}

// SwitchToOrRotateAllBranchesLogs is the universal entry point for the
// "show all-branches log graph" action. Switches to the all-branches
// view, or, if already on it, rotates to the next configured command and
// renders it. Body lifted from
// StatusController.switchToOrRotateAllBranchesLogs.
func (self *AllBranchesLogHelper) SwitchToOrRotateAllBranchesLogs() error {
	// A bit of a hack to ensure we only rotate to the next branch log command
	// if we currently are looking at a branch log. Otherwise, we should just show
	// the current index (if we are coming from the dashboard).
	if self.c.Views().Main.Title != self.c.Tr.StatusTitle {
		self.c.Git().Branch.RotateAllBranchesLogIdx()
	}
	self.ShowAllBranchLogs()
	return nil
}

// SwitchToOrRotateAllBranchesLogsBackward is the reverse direction.
func (self *AllBranchesLogHelper) SwitchToOrRotateAllBranchesLogsBackward() error {
	// A bit of a hack to ensure we only rotate to the previous branch log command
	// if we currently are looking at a branch log. Otherwise, we should just show
	// the current index (if we are coming from the dashboard).
	if self.c.Views().Main.Title != self.c.Tr.StatusTitle {
		self.c.Git().Branch.RotateAllBranchesLogIdxBackward()
	}
	self.ShowAllBranchLogs()
	return nil
}
