package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

type checkedOutRefGuiCommon struct {
	types.IGuiCommon
	model *types.Model
}

func (self *checkedOutRefGuiCommon) Model() *types.Model {
	return self.model
}

func TestFindCheckedOutRef(t *testing.T) {
	testCases := []struct {
		name           string
		branches       []*models.Branch
		expectedBranch *models.Branch
		expectedIndex  int
		expectedFound  bool
	}{
		{name: "empty", expectedIndex: -1},
		{
			name:           "head at zero",
			branches:       []*models.Branch{{Name: "main", Head: true}, {Name: "feature"}},
			expectedBranch: &models.Branch{Name: "main", Head: true},
			expectedIndex:  0,
			expectedFound:  true,
		},
		{
			name:           "nested head",
			branches:       []*models.Branch{{Name: "main"}, {Name: "feature", Head: true}},
			expectedBranch: &models.Branch{Name: "feature", Head: true},
			expectedIndex:  1,
			expectedFound:  true,
		},
		{
			name:          "no head marker",
			branches:      []*models.Branch{{Name: "main"}},
			expectedIndex: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			branch, index, found := findCheckedOutRef(testCase.branches)
			assert.Equal(t, testCase.expectedBranch, branch)
			assert.Equal(t, testCase.expectedIndex, index)
			assert.Equal(t, testCase.expectedFound, found)
		})
	}
}

func TestCheckedOutBranchActionsReturnErrorWithoutHeadMarker(t *testing.T) {
	const noBranchesError = "no branches for this repo"

	guiCommon := &checkedOutRefGuiCommon{
		model: &types.Model{Branches: []*models.Branch{{Name: "main"}}},
	}
	helperCommon := &HelperCommon{
		Common:     &common.Common{Tr: &i18n.TranslationSet{NoBranchesThisRepo: noBranchesError}},
		IGuiCommon: guiCommon,
	}
	refsHelper := NewRefsHelper(helperCommon, nil)
	mergeAndRebaseHelper := NewMergeAndRebaseHelper(helperCommon)

	testCases := []struct {
		name   string
		action func() error
	}{
		{name: "move commits", action: refsHelper.MoveCommitsToNewBranch},
		{name: "rebase", action: func() error { return mergeAndRebaseHelper.RebaseOntoRef("feature") }},
		{name: "merge", action: func() error { return mergeAndRebaseHelper.MergeRefIntoCheckedOutBranch("feature") }},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var err error
			assert.NotPanics(t, func() {
				err = testCase.action()
			})
			assert.EqualError(t, err, noBranchesError)
		})
	}
}
