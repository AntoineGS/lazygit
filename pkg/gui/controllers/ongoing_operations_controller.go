package controllers

import (
	"fmt"
	"sync"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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

// completedOp is a snapshot of an OngoingOperation taken at the moment it was
// removed from the registry, so the popup can keep showing it in a "completed"
// section for the rest of the popup-open session.
type completedOp struct {
	ID        int64
	Label     string
	StartTime time.Time
	EndTime   time.Time
	LastCmd   string
}

// popupSession is the controller's state for a single popup-open. It tracks
// which ops are currently live in the registry and which have completed
// while this popup was open. renderedIDs is parallel to the most-recently
// rendered menu items list; it lets us preserve selection across re-renders
// by mapping selected-index → op-ID before rebuilding, then ID → new-index
// after.
type popupSession struct {
	activeOps   map[int64]*helpers.OngoingOperation
	completed   []completedOp
	renderedIDs []int64
}

func (self *OngoingOperationsController) show() error {
	session := &popupSession{
		activeOps: make(map[int64]*helpers.OngoingOperation),
	}
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
		Items:           self.buildItems(session, onClose),
		HideCancel:      true,
		HideConfirmHint: true,
		OnCancel:        onClose,
	}); err != nil {
		closePopup()
		return err
	}

	go self.refreshOnEvents(session, events, stop, onClose)
	return nil
}

// refreshOnEvents re-renders the popup whenever the registry signals a change.
// Active and completed ops are kept in a session-local list; rows aren't
// removed when an op completes, only restyled — so the user can see what just
// finished. Selection is preserved across re-renders by op-ID, so the
// highlighted row follows the same operation even if its position shifts.
func (self *OngoingOperationsController) refreshOnEvents(session *popupSession, events <-chan struct{}, stop <-chan struct{}, onClose func() error) {
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

				selectedID := self.currentSelectedID(session)
				items := self.buildItems(session, onClose)
				self.c.Contexts().Menu.SetMenuItems(items, nil)
				self.restoreSelection(session, selectedID)
				self.c.PostRefreshUpdate(self.c.Contexts().Menu)
				return nil
			})
		}
	}
}

// currentSelectedID returns the op-ID of the currently-highlighted row, or 0
// if nothing matches (empty list, selection out of bounds, or sentinel row).
func (self *OngoingOperationsController) currentSelectedID(session *popupSession) int64 {
	idx := self.c.Contexts().Menu.GetSelectedLineIdx()
	if idx < 0 || idx >= len(session.renderedIDs) {
		return 0
	}
	return session.renderedIDs[idx]
}

// restoreSelection moves the cursor back to whichever row carries the given
// op-ID after a re-render. If the ID is no longer in the list, the cursor
// stays where it is.
func (self *OngoingOperationsController) restoreSelection(session *popupSession, id int64) {
	if id == 0 {
		return
	}
	for i, otherID := range session.renderedIDs {
		if otherID == id {
			self.c.Contexts().Menu.SetSelection(i)
			return
		}
	}
}

// buildItems re-reads the registry, snapshots any newly-completed ops into
// the session, and renders the combined active + completed list. The
// session.renderedIDs slice is updated in lockstep so selection can be
// preserved by ID across re-renders.
//
// LabelColumns is set explicitly (rather than letting Label propagate via
// createMenu) so that event-driven re-renders — which call SetMenuItems
// directly — produce visible rows.
//
// Each item's OnPress invokes onClose so pressing Enter dismisses the popup
// just like Esc does (and tears down the subscription/goroutine).
func (self *OngoingOperationsController) buildItems(session *popupSession, onClose func() error) []*types.MenuItem {
	currentActive := self.c.Helpers().OngoingOperations.List()

	// Snapshot ops that left the registry since the last render.
	stillActive := make(map[int64]struct{}, len(currentActive))
	for _, op := range currentActive {
		stillActive[op.ID] = struct{}{}
	}
	for id, op := range session.activeOps {
		if _, alive := stillActive[id]; alive {
			continue
		}
		session.completed = append(session.completed, completedOp{
			ID:        op.ID,
			Label:     op.Label,
			StartTime: op.StartTime,
			EndTime:   time.Now(),
			LastCmd:   op.LastCommand(),
		})
	}

	// Refresh the active map.
	session.activeOps = make(map[int64]*helpers.OngoingOperation, len(currentActive))
	for _, op := range currentActive {
		session.activeOps[op.ID] = op
	}

	if len(currentActive) == 0 && len(session.completed) == 0 {
		session.renderedIDs = []int64{0}
		return []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.NoOngoingOperations},
				OnPress:      onClose,
			},
		}
	}

	items := make([]*types.MenuItem, 0, len(currentActive)+len(session.completed))
	ids := make([]int64, 0, len(currentActive)+len(session.completed))

	// Latest at the top, oldest at the bottom: registry returns active ops in
	// StartTime-ascending order and we append completed ops in completion
	// order, so we iterate both backwards.
	for i := len(currentActive) - 1; i >= 0; i-- {
		op := currentActive[i]
		items = append(items, self.buildActiveItem(op, onClose))
		ids = append(ids, op.ID)
	}
	for i := len(session.completed) - 1; i >= 0; i-- {
		c := session.completed[i]
		items = append(items, self.buildCompletedItem(c, onClose))
		ids = append(ids, c.ID)
	}

	session.renderedIDs = ids
	return items
}

func (self *OngoingOperationsController) buildActiveItem(op *helpers.OngoingOperation, onClose func() error) *types.MenuItem {
	cmd := op.CurrentCommand()
	if cmd == "" {
		cmd = "—"
	}
	label := fmt.Sprintf(self.c.Tr.OngoingOperationLineFormat, op.Label, cmd)
	tooltip := fmt.Sprintf(self.c.Tr.OngoingOperationStartedAtFormat, op.StartTime.Format(time.TimeOnly))
	return &types.MenuItem{
		LabelColumns: []string{label},
		Tooltip:      tooltip,
		OnPress:      onClose,
	}
}

func (self *OngoingOperationsController) buildCompletedItem(c completedOp, onClose func() error) *types.MenuItem {
	cmd := c.LastCmd
	if cmd == "" {
		cmd = "—"
	}
	label := fmt.Sprintf(self.c.Tr.OngoingOperationLineFormat, c.Label, cmd)
	tooltip := fmt.Sprintf(self.c.Tr.OngoingOperationCompletedAtFormat,
		c.StartTime.Format(time.TimeOnly), c.EndTime.Format(time.TimeOnly))
	// Trailing unstyled space: gocui reuses the last character's fgColor for
	// empty cells past end-of-line (vendor view.go:1288), so without it the
	// strikethrough bleeds across the row when the line is highlighted.
	styled := style.FgDefault.SetStrikethrough().Sprint(label) + " "
	return &types.MenuItem{
		LabelColumns: []string{styled},
		Tooltip:      tooltip,
		OnPress:      onClose,
	}
}
