package controllers

import (
	"fmt"
	"time"

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
	if err := self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.OngoingOperations,
		Items:           self.buildItems(),
		HideCancel:      true,
		HideConfirmHint: true,
	}); err != nil {
		return err
	}

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
			self.c.Contexts().Menu.SetMenuItems(self.buildItems(), nil)
			self.c.PostRefreshUpdate(self.c.Contexts().Menu)
			return nil
		})
	}
}

// buildItems renders the current registry as menu items. Each operation is one
// row; when nothing is running, a single inert row carries the empty-state
// message so the popup is always a list (avoids prompt-only rendering quirks
// like a stale Tooltip pane appearing under the menu).
//
// LabelColumns is set explicitly (rather than letting Label propagate via
// createMenu) so that auto-refresh — which calls SetMenuItems directly without
// going back through createMenu — produces visible rows.
func (self *OngoingOperationsController) buildItems() []*types.MenuItem {
	noop := func() error { return nil }

	ops := self.c.Helpers().OngoingOperations.List()
	if len(ops) == 0 {
		return []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.NoOngoingOperations},
				OnPress:      noop,
			},
		}
	}

	items := make([]*types.MenuItem, 0, len(ops))
	for _, op := range ops {
		duration := op.Elapsed().Truncate(time.Second)
		cmd := op.CurrentCommand()
		if cmd == "" {
			cmd = "—"
		}
		label := fmt.Sprintf(self.c.Tr.OngoingOperationLineFormat, op.Label, duration.String(), cmd)
		items = append(items, &types.MenuItem{
			LabelColumns: []string{label},
			OnPress:      noop,
		})
	}
	return items
}
