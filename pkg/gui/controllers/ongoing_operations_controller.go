package controllers

import (
	"fmt"
	"sync"
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

func (self *OngoingOperationsController) show() error {
	events, unsubscribe := self.c.Helpers().OngoingOperations.Subscribe()

	stop := make(chan struct{})
	closeOnce := sync.Once{}
	closePopup := func() {
		closeOnce.Do(func() {
			close(stop)
			unsubscribe()
		})
	}
	onClose := func() error {
		closePopup()
		return nil
	}

	if err := self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.OngoingOperations,
		Items:           self.buildItems(onClose),
		HideCancel:      true,
		HideConfirmHint: true,
		OnCancel:        onClose,
	}); err != nil {
		closePopup()
		return err
	}

	go self.refreshOnEvents(events, stop, onClose)
	return nil
}

// refreshOnEvents re-renders the popup whenever the registry signals a change.
// Exits when stop is closed (popup dismissed). Durations don't auto-tick
// between events — that's by design for the "is this stuck?" use case (a
// frozen duration means nothing has happened).
func (self *OngoingOperationsController) refreshOnEvents(events <-chan struct{}, stop <-chan struct{}, onClose func() error) {
	for {
		select {
		case <-stop:
			return
		case <-events:
			self.c.OnUIThread(func() error {
				select {
				case <-stop:
					return nil
				default:
				}
				self.c.Contexts().Menu.SetMenuItems(self.buildItems(onClose), nil)
				self.c.PostRefreshUpdate(self.c.Contexts().Menu)
				return nil
			})
		}
	}
}

// buildItems renders the current registry as menu items. Each operation is one
// row; when nothing is running, a single inert row carries the empty-state
// message so the popup is always a list (avoids prompt-only rendering quirks
// like a stale Tooltip pane appearing under the menu).
//
// LabelColumns is set explicitly (rather than letting Label propagate via
// createMenu) so that the event-driven re-render — which calls SetMenuItems
// directly without going back through createMenu — produces visible rows.
//
// Each item's OnPress invokes onClose so pressing Enter dismisses the popup
// just like Esc does (and also tears down the subscription/goroutine).
func (self *OngoingOperationsController) buildItems(onClose func() error) []*types.MenuItem {
	ops := self.c.Helpers().OngoingOperations.List()
	if len(ops) == 0 {
		return []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.NoOngoingOperations},
				OnPress:      onClose,
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
			OnPress:      onClose,
		})
	}
	return items
}
