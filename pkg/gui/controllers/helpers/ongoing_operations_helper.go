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
	// atomically without locking the registry mutex.
	currentCommand atomic.Pointer[string]
}

// SetCurrentCommand records the git command currently being executed inside
// this operation. Pass "" to indicate no command is currently running.
func (op *OngoingOperation) SetCurrentCommand(cmd string) {
	op.currentCommand.Store(&cmd)
}

// CurrentCommand returns the most recently recorded command, or "".
func (op *OngoingOperation) CurrentCommand() string {
	if p := op.currentCommand.Load(); p != nil {
		return *p
	}
	return ""
}

// Elapsed returns the wall-clock duration since registration.
func (op *OngoingOperation) Elapsed() time.Duration {
	return time.Since(op.StartTime)
}

type OngoingOperationsHelper struct {
	mutex      deadlock.Mutex
	operations map[int64]*OngoingOperation
	nextID     int64
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
	}
	self.operations[op.ID] = op
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
