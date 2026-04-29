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

// refreshInterval is how often the popup re-renders while open, so newly-started
// or completed operations appear/disappear without the user re-pressing the key.
const refreshInterval = time.Second

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

	go self.refreshWhileOpen()
	return nil
}

// refreshWhileOpen periodically re-renders the popup so operations that start
// or finish while it's open appear/disappear in place. Exits when the user
// dismisses the popup (the menu is no longer the current context).
func (self *OngoingOperationsController) refreshWhileOpen() {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	menuKey := self.c.Contexts().Menu.GetKey()
	for range ticker.C {
		if self.c.Context().Current().GetKey() != menuKey {
			return
		}
		self.c.OnUIThread(func() error {
			// Re-check inside the UI thread: the user may have dismissed the
			// popup, or another menu may have replaced it, between the ticker
			// firing and this callback running.
			if self.c.Context().Current().GetKey() != menuKey {
				return nil
			}
			ops := self.c.Helpers().OngoingOperations.List()
			self.c.Contexts().Menu.SetPrompt(self.formatOps(ops))
			self.c.PostRefreshUpdate(self.c.Contexts().Menu)
			return nil
		})
	}
}

func (self *OngoingOperationsController) formatOps(ops []*helpers.OngoingOperation) string {
	if len(ops) == 0 {
		return self.c.Tr.NoOngoingOperations
	}
	var sb strings.Builder
	for _, op := range ops {
		duration := op.Elapsed().Truncate(time.Second)
		cmd := op.CurrentCommand()
		if cmd == "" {
			cmd = "—"
		}
		fmt.Fprintf(&sb, self.c.Tr.OngoingOperationLineFormat+"\n", op.Label, duration.String(), cmd)
	}
	return strings.TrimRight(sb.String(), "\n")
}
