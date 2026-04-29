package helpers

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOngoingOperationsHelper_RegisterAssignsUniqueIDs(t *testing.T) {
	h := NewOngoingOperationsHelper()

	op1 := h.Register("Pulling 'main'")
	op2 := h.Register("Pushing 'main'")

	assert.NotEqual(t, op1.ID, op2.ID)
	assert.Equal(t, "Pulling 'main'", op1.Label)
	assert.Equal(t, "Pushing 'main'", op2.Label)
	assert.False(t, op1.StartTime.IsZero())
}

func TestOngoingOperationsHelper_ListReturnsActiveOperationsByStartTime(t *testing.T) {
	h := NewOngoingOperationsHelper()

	op1 := h.Register("Pulling 'main'")
	time.Sleep(2 * time.Millisecond)
	op2 := h.Register("Fetching")

	got := h.List()

	assert.Len(t, got, 2)
	assert.Equal(t, op1.ID, got[0].ID)
	assert.Equal(t, op2.ID, got[1].ID)
}

func TestOngoingOperationsHelper_UnregisterRemovesFromList(t *testing.T) {
	h := NewOngoingOperationsHelper()

	op1 := h.Register("Pulling 'main'")
	op2 := h.Register("Fetching")
	h.Unregister(op1)

	got := h.List()

	assert.Len(t, got, 1)
	assert.Equal(t, op2.ID, got[0].ID)
}

func TestOngoingOperationsHelper_UnregisterNilIsNoop(t *testing.T) {
	h := NewOngoingOperationsHelper()

	h.Unregister(nil)

	assert.Empty(t, h.List())
}

func TestOngoingOperation_SetAndGetCurrentCommand(t *testing.T) {
	h := NewOngoingOperationsHelper()
	op := h.Register("Pulling 'main'")

	assert.Equal(t, "", op.CurrentCommand())

	op.SetCurrentCommand("git pull --no-edit origin refs/heads/main")
	assert.Equal(t, "git pull --no-edit origin refs/heads/main", op.CurrentCommand())

	op.SetCurrentCommand("")
	assert.Equal(t, "", op.CurrentCommand())
}

func TestOngoingOperationsHelper_ConcurrentRegisterUnregister(t *testing.T) {
	h := NewOngoingOperationsHelper()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			op := h.Register("op")
			op.SetCurrentCommand("git status")
			h.Unregister(op)
		}()
	}

	wg.Wait()
	assert.Empty(t, h.List())
}

func TestOngoingOperationsHelper_SubscribeReceivesEventsForRegisterUnregisterAndCommandChange(t *testing.T) {
	h := NewOngoingOperationsHelper()
	events, unsubscribe := h.Subscribe()
	defer unsubscribe()

	op := h.Register("Pulling 'main'")
	assertEventReceived(t, events, "Register")

	op.SetCurrentCommand("git pull")
	assertEventReceived(t, events, "SetCurrentCommand non-empty")

	h.Unregister(op)
	assertEventReceived(t, events, "Unregister")
}

func TestOngoingOperationsHelper_SetCurrentCommandClearDoesNotNotify(t *testing.T) {
	// Clearing the command (passing "") must not fire an event: cmd-end is
	// immediately followed by Unregister, and a clear-notify would render the
	// row with cmd="" briefly before removing it.
	h := NewOngoingOperationsHelper()
	op := h.Register("Pulling 'main'")
	op.SetCurrentCommand("git pull")

	events, unsubscribe := h.Subscribe()
	defer unsubscribe()

	op.SetCurrentCommand("")

	select {
	case <-events:
		t.Fatal("expected no event when clearing command; got one")
	case <-time.After(10 * time.Millisecond):
	}
	assert.Equal(t, "", op.CurrentCommand(), "atomic value should still be cleared")
}

func TestOngoingOperationsHelper_SubscribeCoalescesBurstsIntoOneWakeup(t *testing.T) {
	h := NewOngoingOperationsHelper()
	events, unsubscribe := h.Subscribe()
	defer unsubscribe()

	// Three rapid events; the buffered (size-1) channel should hold one
	// pending wake-up at most.
	op1 := h.Register("a")
	op2 := h.Register("b")
	op3 := h.Register("c")
	_, _, _ = op1, op2, op3

	assertEventReceived(t, events, "first")

	// Channel should now be empty — the second and third events were
	// coalesced into the one pending slot.
	select {
	case <-events:
		t.Fatal("expected coalesced events; got a second wake-up")
	case <-time.After(10 * time.Millisecond):
	}
}

func TestOngoingOperationsHelper_UnsubscribeStopsDelivery(t *testing.T) {
	h := NewOngoingOperationsHelper()
	events, unsubscribe := h.Subscribe()

	h.Register("first")
	assertEventReceived(t, events, "before unsubscribe")

	unsubscribe()
	h.Register("second")

	select {
	case <-events:
		t.Fatal("received event after unsubscribe")
	case <-time.After(10 * time.Millisecond):
	}
}

func assertEventReceived(t *testing.T, events <-chan struct{}, label string) {
	t.Helper()
	select {
	case <-events:
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("expected %s event within 50ms; got none", label)
	}
}
