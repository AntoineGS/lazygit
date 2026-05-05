package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestBuildChordContinuations_HiddenBindingsAreFiltered(t *testing.T) {
	keyMn, _ := config.KeyFromLabel("Mn")
	keyMf, _ := config.KeyFromLabel("Mf")
	keyM, _ := config.KeyFromLabel("M")

	visible := &types.Binding{
		Key:         keyMn,
		Description: "Non-FF",
	}
	hidden := &types.Binding{
		Key:                keyMf,
		Description:        "FF only",
		HiddenInChordPopup: func() bool { return true },
	}

	rows := BuildChordContinuations(
		[]*types.Binding{visible, hidden},
		keyM.Sequence(),
		map[string]map[string]config.KeybindingGroupConfig{},
		"localBranches",
	)

	assert.Len(t, rows, 1, "expected only the visible row to remain")
	assert.Equal(t, "Non-FF", rows[0].description)
}
