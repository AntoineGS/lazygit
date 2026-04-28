package gocui

import (
	"testing"
)

func newTestGui(t *testing.T) *Gui {
	t.Helper()
	g, err := NewGui(NewGuiOpts{
		OutputMode: OutputNormal,
		Headless:   true,
		Width:      80,
		Height:     24,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { g.Close() })
	return g
}

// Creates a few views and does an initial full flush so all views
// start in a clean (non-tainted) state.
func setupViews(t *testing.T, g *Gui) (*View, *View) {
	t.Helper()

	status, _ := g.SetView("status", 0, 22, 40, 24, 0)
	status.Frame = false
	main, _ := g.SetView("main", 0, 0, 80, 22, 0)

	// Initial content
	status.SetContent("Ready")
	main.SetContent("hello world")

	// Full flush to draw everything and clear tainted flags
	if err := g.flush(); err != nil {
		t.Fatal(err)
	}

	return status, main
}

// Pushes a content-only event directly to the channel
// (synchronous, deterministic — unlike Update which spawns a goroutine).
func pushContentOnly(g *Gui, f func(*Gui) error) {
	g.userEvents <- userEvent{f: f, task: g.NewTask(), contentOnly: true}
}

// Pushes a regular event directly to the channel.
func pushRegular(g *Gui, f func(*Gui) error) {
	g.userEvents <- userEvent{f: f, task: g.NewTask(), contentOnly: false}
}

func TestFlushContentOnly_SkipsUntaintedViews(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// After initial flush, both views should be untainted
	if status.IsTainted() {
		t.Fatal("status view should not be tainted after flush")
	}
	if main.IsTainted() {
		t.Fatal("main view should not be tainted after flush")
	}

	// Modify only the status view
	status.SetContent("Fetching /")

	if !status.IsTainted() {
		t.Fatal("status view should be tainted after SetContent")
	}
	if main.IsTainted() {
		t.Fatal("main view should not be tainted (was not modified)")
	}

	// flushContentOnly should succeed and clear status tainted flag
	if err := g.flushContentOnly(); err != nil {
		t.Fatal(err)
	}

	if status.IsTainted() {
		t.Fatal("status view should not be tainted after flushContentOnly")
	}
	if main.IsTainted() {
		t.Fatal("main view should not be tainted after flushContentOnly")
	}
}

func TestFlushContentOnly_WritesCorrectContent(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	status.SetContent("Fetching |")
	if err := g.flushContentOnly(); err != nil {
		t.Fatal(err)
	}

	buf := status.Buffer()
	if buf != "Fetching |" {
		t.Fatalf("expected status buffer %q, got %q", "Fetching |", buf)
	}
}

func TestProcessEvent_ContentOnlyEvent_SkipsTaintedCheck(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// Send a content-only event that modifies only the status view
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching /")
		return nil
	})

	if err := g.processEvent(); err != nil {
		t.Fatal(err)
	}

	if status.IsTainted() {
		t.Fatal("status should not be tainted after processEvent with contentOnly")
	}

	if main.IsTainted() {
		t.Fatal("main should not be tainted since it was not modified")
	}
}

func TestProcessEvent_RegularEvent_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	// Regular event (not content-only) should trigger full flush
	pushRegular(g, func(gui *Gui) error {
		status.SetContent("Fetching \\")
		return nil
	})

	if err := g.processEvent(); err != nil {
		t.Fatal(err)
	}

	if status.IsTainted() {
		t.Fatal("status should not be tainted after full flush")
	}
}

// Queue a content-only event followed by a regular event.
// processEvent picks up the first; processRemainingEvents picks up
// the second. Since the second is not contentOnly, full flush runs.
func TestProcessEvent_MixedBatch_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching -")
		return nil
	})
	pushRegular(g, func(gui *Gui) error {
		main.SetContent("updated main")
		return nil
	})

	if err := g.processEvent(); err != nil {
		t.Fatal(err)
	}

	// Both views were modified and should have been drawn by full flush
	if status.IsTainted() {
		t.Fatal("status should not be tainted after full flush")
	}
	if main.IsTainted() {
		t.Fatal("main should not be tainted after full flush")
	}
}

// Even if a regular event comes first and the remaining are contentOnly,
// the batch must use full flush.
func TestProcessEvent_RegularThenContentOnly_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	pushRegular(g, func(gui *Gui) error {
		main.SetContent("new main content")
		return nil
	})
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching |")
		return nil
	})

	if err := g.processEvent(); err != nil {
		t.Fatal(err)
	}

	if status.IsTainted() {
		t.Fatal("status should not be tainted after full flush")
	}
	if main.IsTainted() {
		t.Fatal("main should not be tainted after full flush")
	}
}

func TestProcessRemainingEvents_AllContentOnly_ReturnsTrue(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("a")
		return nil
	})
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("b")
		return nil
	})

	contentOnly, err := g.processRemainingEvents()
	if err != nil {
		t.Fatal(err)
	}
	if !contentOnly {
		t.Fatal("should return true when all events are contentOnly")
	}
}

func TestProcessRemainingEvents_MixedEvents_ReturnsFalse(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("a")
		return nil
	})
	pushRegular(g, func(gui *Gui) error {
		status.SetContent("b")
		return nil
	})

	contentOnly, err := g.processRemainingEvents()
	if err != nil {
		t.Fatal(err)
	}
	if contentOnly {
		t.Fatal("should return false when any event is not contentOnly")
	}
}

func TestProcessRemainingEvents_EmptyQueue_ReturnsTrue(t *testing.T) {
	g := newTestGui(t)

	contentOnly, err := g.processRemainingEvents()
	if err != nil {
		t.Fatal(err)
	}
	if !contentOnly {
		t.Fatal("should return true when no events are queued")
	}
}
