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
