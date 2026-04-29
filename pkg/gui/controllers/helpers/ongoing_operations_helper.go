package helpers

import (
	"slices"
	"sync/atomic"
	"time"

	"github.com/sasha-s/go-deadlock"
)

// OngoingOperation represents a long-running, user-visible operation
// (pull, push, fetch, checkout, etc.). It is registered on entry to
// WithInlineStatus / WithWaitingStatusImpl and unregistered on exit.
type OngoingOperation struct {
	ID        int64
	Label     string
	StartTime time.Time

	// currentCommand holds a *string so SetCurrentCommand can swap
	// atomically without locking the registry mutex. It is cleared
	// (set to "") at cmd-end.
	currentCommand atomic.Pointer[string]

	// lastCommand mirrors the most recent non-empty command. It is never
	// cleared, so a popup snapshotting an operation that has already
	// completed (cmd cleared, op unregistered) can still surface what ran.
	lastCommand atomic.Pointer[string]

	// notify is set by Register so SetCurrentCommand can wake any popup
	// subscribers without taking a back-reference to the helper.
	notify func()
}

// SetCurrentCommand records the git command currently being executed inside
// this operation. Pass "" to indicate no command is currently running.
//
// We only notify subscribers when a non-empty command is set, not when one
// is cleared. The cmd runner clears at cmd-end, immediately followed by the
// containing operation's Unregister; firing both events would briefly render
// the row with cmd="" before removing it. Skipping the clear-notify lets the
// next event (Unregister, or a follow-on SetCurrentCommand for multi-cmd ops)
// pick up the cleared atomic value naturally without an intermediate flash.
//
// Non-empty values are also recorded as lastCommand, which never gets
// cleared — see LastCommand.
func (op *OngoingOperation) SetCurrentCommand(cmd string) {
	op.currentCommand.Store(&cmd)
	if cmd == "" {
		return
	}
	op.lastCommand.Store(&cmd)
	if op.notify != nil {
		op.notify()
	}
}

// CurrentCommand returns the command currently executing, or "" between
// commands or when none has run yet.
func (op *OngoingOperation) CurrentCommand() string {
	if p := op.currentCommand.Load(); p != nil {
		return *p
	}
	return ""
}

// LastCommand returns the most recent non-empty command. Unlike
// CurrentCommand, it is not cleared when a command finishes — so it remains
// usable after the operation has been unregistered, for popups that want to
// keep showing what ran.
func (op *OngoingOperation) LastCommand() string {
	if p := op.lastCommand.Load(); p != nil {
		return *p
	}
	return ""
}

// Elapsed returns the wall-clock duration since registration.
func (op *OngoingOperation) Elapsed() time.Duration {
	return time.Since(op.StartTime)
}

type OngoingOperationsHelper struct {
	mutex       deadlock.Mutex
	operations  map[int64]*OngoingOperation
	nextID      int64
	subscribers []chan struct{}
}

// NewOngoingOperationsHelper returns an empty registry. The returned helper is
// safe to call from multiple goroutines.
func NewOngoingOperationsHelper() *OngoingOperationsHelper {
	return &OngoingOperationsHelper{
		operations: make(map[int64]*OngoingOperation),
	}
}

// Register adds a new operation to the registry and returns its handle. The
// caller must call Unregister with the same handle when the operation ends.
func (self *OngoingOperationsHelper) Register(label string) *OngoingOperation {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.nextID++
	op := &OngoingOperation{
		ID:        self.nextID,
		Label:     label,
		StartTime: time.Now(),
		notify:    self.Notify,
	}
	self.operations[op.ID] = op
	self.notifyLocked()
	return op
}

// Unregister removes an operation from the registry. Passing nil is a no-op,
// to make defer call sites simple.
func (self *OngoingOperationsHelper) Unregister(op *OngoingOperation) {
	if op == nil {
		return
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	delete(self.operations, op.ID)
	self.notifyLocked()
}

// List returns a snapshot of active operations sorted by start time
// (oldest first), so the popup shows the most-likely-stuck ones at the top.
func (self *OngoingOperationsHelper) List() []*OngoingOperation {
	self.mutex.Lock()
	ops := make([]*OngoingOperation, 0, len(self.operations))
	for _, op := range self.operations {
		ops = append(ops, op)
	}
	self.mutex.Unlock()

	slices.SortFunc(ops, func(a, b *OngoingOperation) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return ops
}

// Subscribe returns a channel that receives a signal each time the registry
// changes (Register, Unregister, or SetCurrentCommand on any operation). The
// channel is buffered with a single slot, so emitters never block — bursts
// of events while a subscriber is busy are coalesced into one wake-up.
//
// Callers MUST invoke the returned unsubscribe func when they're done, or
// the helper will retain the channel and emit to it forever.
func (self *OngoingOperationsHelper) Subscribe() (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1)

	self.mutex.Lock()
	self.subscribers = append(self.subscribers, ch)
	self.mutex.Unlock()

	unsubscribe := func() {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		self.subscribers = slices.DeleteFunc(self.subscribers, func(c chan struct{}) bool {
			return c == ch
		})
	}
	return ch, unsubscribe
}

// Notify wakes all subscribers. Safe to call without holding the registry
// mutex; used by OngoingOperation.SetCurrentCommand.
func (self *OngoingOperationsHelper) Notify() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.notifyLocked()
}

// notifyLocked broadcasts to subscribers; the registry mutex must be held.
func (self *OngoingOperationsHelper) notifyLocked() {
	for _, ch := range self.subscribers {
		select {
		case ch <- struct{}{}:
		default:
			// channel already has a pending wake-up; coalesce
		}
	}
}
