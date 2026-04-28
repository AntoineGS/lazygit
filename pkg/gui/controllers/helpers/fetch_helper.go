package helpers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FetchHelper struct {
	c              *HelperCommon
	branchesHelper *BranchesHelper
}

func NewFetchHelper(c *HelperCommon, branchesHelper *BranchesHelper) *FetchHelper {
	return &FetchHelper{
		c:              c,
		branchesHelper: branchesHelper,
	}
}

// Fetch runs a fetch on the current repo. Body lifted from
// FilesController.fetch.
func (self *FetchHelper) Fetch() error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingStatus, func(task gocui.Task) error {
		self.c.LogAction("Fetch")
		err := self.c.Git().Sync.Fetch(task)

		if err != nil && strings.Contains(err.Error(), "exit status 128") {
			return errors.New(self.c.Tr.PassUnameWrong)
		}

		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.SYNC})

		if err == nil {
			err = self.branchesHelper.AutoForwardBranches()
		}

		return err
	})
}
