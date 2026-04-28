package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type OngoingOperationsController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &OngoingOperationsController{}

func NewOngoingOperationsController(c *ControllerCommon) *OngoingOperationsController {
	return &OngoingOperationsController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *OngoingOperationsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.ShowOngoingOperations),
			Handler:     opts.Guards.NoPopupPanel(self.show),
			Description: self.c.Tr.OngoingOperations,
			Tooltip:     self.c.Tr.OngoingOperationsTooltip,
			OpensMenu:   true,
		},
	}
}

func (self *OngoingOperationsController) Context() types.Context {
	return nil
}

func (self *OngoingOperationsController) show() error {
	ops := self.c.Helpers().OngoingOperations.List()

	if err := self.c.Menu(types.CreateMenuOptions{
		Title:      self.c.Tr.OngoingOperations,
		Prompt:     self.formatOps(ops),
		Items:      []*types.MenuItem{},
		HideCancel: true,
	}); err != nil {
		return err
	}

	// MenuController only sets the tooltip when an item is selected; with no
	// items the buffer would otherwise hold stale content from the previous
	// menu, leaving an empty tooltip pane visible below the prompt.
	self.c.Views().Tooltip.SetContent("")
	return nil
}

func (self *OngoingOperationsController) formatOps(ops []*helpers.OngoingOperation) string {
	if len(ops) == 0 {
		return self.c.Tr.NoOngoingOperations
	}
	var sb strings.Builder
	now := time.Now()
	for _, op := range ops {
		duration := now.Sub(op.StartTime).Truncate(time.Second)
		cmd := op.CurrentCommand()
		if cmd == "" {
			cmd = "—"
		}
		fmt.Fprintf(&sb, self.c.Tr.OngoingOperationLineFormat+"\n", op.Label, duration.String(), cmd)
	}
	return strings.TrimRight(sb.String(), "\n")
}
