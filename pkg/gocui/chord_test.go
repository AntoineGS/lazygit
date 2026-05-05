package gocui

import (
	"errors"
	"testing"
)

func TestAllowChordStartsCallbackBlocksNewChord(t *testing.T) {
	g := &Gui{}
	view := NewView("commitMessage", 0, 0, 10, 10, OutputNormal)
	g.SetAllowChordStartsCallback(func(*View) bool { return false })
	g.keybindings = []*keybinding{
		newChordKeybinding(view.Name(), []Key{NewKeyRune('g'), NewKeyRune('p')}, func(*Gui, *View) error {
			return nil
		}),
	}

	if err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('g')}); err != nil {
		t.Fatalf("execKeybindings returned error: %v", err)
	}
	if got := g.PendingChord(); got != nil {
		t.Fatalf("expected no pending chord when starts are disabled, got %#v", got)
	}
}

func TestAllowChordStartsCallbackAllowsNewChord(t *testing.T) {
	g := &Gui{}
	view := NewView("files", 0, 0, 10, 10, OutputNormal)
	g.SetAllowChordStartsCallback(func(*View) bool { return true })
	g.keybindings = []*keybinding{
		newChordKeybinding(view.Name(), []Key{NewKeyRune('g'), NewKeyRune('p')}, func(*Gui, *View) error {
			return nil
		}),
	}

	if err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('g')}); err != nil {
		t.Fatalf("execKeybindings returned error: %v", err)
	}
	got := g.PendingChord()
	if len(got) != 1 || !got[0].Equals(NewKeyRune('g')) {
		t.Fatalf("expected pending chord [g], got %#v", got)
	}
}

func TestSetPendingChord_FiresOnChordStateChange(t *testing.T) {
	g := &Gui{}
	var got []Key
	g.SetChordStateCallback(func(prefix []Key) {
		got = append([]Key(nil), prefix...)
	})

	prefix := []Key{NewKeyRune('D')}
	g.SetPendingChord(prefix, "files")

	if len(got) != 1 {
		t.Fatalf("expected onChordStateChange to fire with 1-key prefix, got %d", len(got))
	}
	if g.PendingChord() == nil || len(g.PendingChord()) != 1 {
		t.Fatalf("expected pendingChord to have 1 key, got %v", g.PendingChord())
	}
	if g.PendingChordView() != "files" {
		t.Fatalf("expected pendingChordView to be %q, got %q", "files", g.PendingChordView())
	}
}

func TestSetPendingChord_EmptyPrefixRoutesThroughClear(t *testing.T) {
	g := &Gui{}
	calls := 0
	var lastPrefix []Key
	g.SetChordStateCallback(func(p []Key) { calls++; lastPrefix = p })

	g.SetPendingChord([]Key{NewKeyRune('g')}, "files")
	if calls != 1 || len(lastPrefix) != 1 {
		t.Fatalf("after first SetPendingChord: calls=%d, prefix=%v", calls, lastPrefix)
	}

	g.SetPendingChord(nil, "files")
	if calls != 2 {
		t.Fatalf("expected exactly one additional callback (clear), got total=%d", calls)
	}
	if len(lastPrefix) != 0 {
		t.Fatalf("expected empty prefix on clear, got %v", lastPrefix)
	}
	if g.pendingChord != nil || g.pendingChordView != "" {
		t.Fatal("state should be cleared")
	}
}

func TestSetPendingChord_MutatesIndependently(t *testing.T) {
	g := &Gui{}
	prefix := []Key{NewKeyRune('g'), NewKeyRune('p')}
	g.SetPendingChord(prefix, "files")

	prefix[0] = NewKeyRune('x')

	got := g.PendingChord()
	if got[0].Equals(NewKeyRune('x')) {
		t.Fatal("SetPendingChord should copy the prefix, not store the original slice")
	}
}

func TestPasteStartClearsPendingChord(t *testing.T) {
	g := &Gui{}
	cleared := false
	g.SetChordStateCallback(func(p []Key) {
		if len(p) == 0 {
			cleared = true
		}
	})
	g.pendingChord = []Key{NewKeyRune('g')}
	g.pendingChordView = "files"

	if err := g.handleEvent(&GocuiEvent{Type: eventPaste, Start: true}); err != nil {
		t.Fatalf("handleEvent(eventPaste): %v", err)
	}

	if !g.IsPasting {
		t.Fatal("IsPasting should be true after paste-start")
	}
	if g.pendingChord != nil || g.pendingChordView != "" {
		t.Fatal("pending chord should be cleared on paste-start")
	}
	if !cleared {
		t.Fatal("expected onChordStateChange to fire with empty prefix")
	}
}

func noopChordHandler(*Gui, *View) error { return nil }

func TestSetKeybindingKeys_RejectsPrefixOfExistingChord(t *testing.T) {
	g := &Gui{}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p'), NewKeyRune('s')}, noopChordHandler); err != nil {
		t.Fatalf("first registration should succeed, got %v", err)
	}
	err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler)
	if err == nil {
		t.Fatal("expected error when registering chord that is a prefix of an existing chord")
	}
}

func TestSetKeybindingKeys_RejectsExtensionOfExistingChord(t *testing.T) {
	g := &Gui{}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
		t.Fatalf("first registration should succeed, got %v", err)
	}
	err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p'), NewKeyRune('s')}, noopChordHandler)
	if err == nil {
		t.Fatal("expected error when registering chord that extends an existing chord")
	}
}

func TestSetKeybindingKeys_AllowsDuplicateForRuntimeMutex(t *testing.T) {
	// lazygit registers two bindings with the same chord on the same
	// view when they are mode-exclusive at runtime via
	// GetDisabledReason / AllowFurtherDispatching (e.g. bisect mid- vs
	// pre-start bindings). Validation must permit that pattern.
	g := &Gui{}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
		t.Fatalf("first registration should succeed, got %v", err)
	}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
		t.Fatalf("second registration of same chord on same view should succeed (runtime mutex), got %v", err)
	}
}

func TestSetKeybindingKeys_AllowsDistinctChords(t *testing.T) {
	g := &Gui{}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
		t.Fatalf("first registration should succeed, got %v", err)
	}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('q')}, noopChordHandler); err != nil {
		t.Fatalf("second registration with non-conflicting chord should succeed, got %v", err)
	}
}

func TestSetKeybindingKeys_DetectsGlobalVsViewConflict(t *testing.T) {
	t.Run("explicit view first, global second", func(t *testing.T) {
		g := &Gui{}
		if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p'), NewKeyRune('s')}, noopChordHandler); err != nil {
			t.Fatalf("first registration should succeed, got %v", err)
		}
		err := g.SetKeybindingKeys("", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler)
		if err == nil {
			t.Fatal("expected error: global chord conflicts with explicit-view chord (prefix)")
		}
	})

	t.Run("global first, explicit view second", func(t *testing.T) {
		g := &Gui{}
		if err := g.SetKeybindingKeys("", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
			t.Fatalf("first registration should succeed, got %v", err)
		}
		err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p'), NewKeyRune('s')}, noopChordHandler)
		if err == nil {
			t.Fatal("expected error: explicit-view chord conflicts with existing global chord (prefix)")
		}
	})
}

func TestSetKeybindingKeys_AllowsSamePrefixDifferentExplicitViews(t *testing.T) {
	g := &Gui{}
	if err := g.SetKeybindingKeys("files", []Key{NewKeyRune('b'), NewKeyRune('p'), NewKeyRune('s')}, noopChordHandler); err != nil {
		t.Fatalf("first registration should succeed, got %v", err)
	}
	if err := g.SetKeybindingKeys("branches", []Key{NewKeyRune('b'), NewKeyRune('p')}, noopChordHandler); err != nil {
		t.Fatalf("registration on different explicit view should succeed, got %v", err)
	}
}

// Lazygit registers two chord bindings with the same chord on the same view
// and disambiguates at runtime by having the inactive binding return
// ErrKeybindingNotHandled (via GetDisabledReason / AllowFurtherDispatching).
// The chord dispatcher must iterate to the next matching binding when the
// first declines.
func TestChordDispatch_FallsThroughOnErrKeybindingNotHandled(t *testing.T) {
	g := &Gui{}
	view := NewView("commits", 0, 0, 10, 10, OutputNormal)
	g.views = append(g.views, view)

	secondFired := false
	g.keybindings = []*keybinding{
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			return ErrKeybindingNotHandled
		}),
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			secondFired = true
			return nil
		}),
	}

	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "commits"

	if err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('b')}); err != nil {
		t.Fatalf("execKeybindings returned error: %v", err)
	}
	if !secondFired {
		t.Fatal("expected second handler to fire after first returned ErrKeybindingNotHandled")
	}
	if g.pendingChord != nil {
		t.Fatalf("expected pendingChord to be cleared, got %v", g.pendingChord)
	}
}

// When the first matching chord handler accepts the binding, the dispatcher
// must NOT fall through to subsequent matching bindings.
func TestChordDispatch_FirstHandlerWinsWhenItAccepts(t *testing.T) {
	g := &Gui{}
	view := NewView("commits", 0, 0, 10, 10, OutputNormal)
	g.views = append(g.views, view)

	flag1 := false
	flag2 := false
	g.keybindings = []*keybinding{
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			flag1 = true
			return nil
		}),
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			flag2 = true
			return nil
		}),
	}

	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "commits"

	if err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('b')}); err != nil {
		t.Fatalf("execKeybindings returned error: %v", err)
	}
	if !flag1 {
		t.Fatal("expected first handler to fire")
	}
	if flag2 {
		t.Fatal("expected second handler NOT to fire when first accepted")
	}
	if g.pendingChord != nil {
		t.Fatalf("expected pendingChord to be cleared, got %v", g.pendingChord)
	}
}

// If every binding for an exact-matched chord declines via
// ErrKeybindingNotHandled, the chord dispatcher should treat it as
// "matched but suppressed": clear the pending chord and return nil.
// It must NOT fall through to the prefix-extension loop (which would
// re-match the just-completed chord as a prefix and leave pendingChord
// stuck on the full sequence).
func TestChordDispatch_AllHandlersDeclineDropsChord(t *testing.T) {
	g := &Gui{}
	view := NewView("commits", 0, 0, 10, 10, OutputNormal)
	g.views = append(g.views, view)

	g.keybindings = []*keybinding{
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			return ErrKeybindingNotHandled
		}),
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			return ErrKeybindingNotHandled
		}),
	}

	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "commits"

	err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('b')})
	if err != nil {
		t.Fatalf("expected nil error when all handlers decline (matched-but-suppressed), got %v", err)
	}
	if g.pendingChord != nil {
		t.Fatalf("expected pendingChord to be cleared, got %v", g.pendingChord)
	}
}

// Sanity check: errors that are NOT ErrKeybindingNotHandled propagate
// out of the dispatcher and stop iteration.
func TestChordDispatch_NonNotHandledErrorStopsIteration(t *testing.T) {
	g := &Gui{}
	view := NewView("commits", 0, 0, 10, 10, OutputNormal)
	g.views = append(g.views, view)

	sentinel := errors.New("boom")
	secondFired := false
	g.keybindings = []*keybinding{
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			return sentinel
		}),
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			secondFired = true
			return nil
		}),
	}

	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "commits"

	err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('b')})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error to propagate, got %v", err)
	}
	if secondFired {
		t.Fatal("expected second handler NOT to fire when first returned non-NotHandled error")
	}
}

// Pins the documented behavior of a CONTRACT VIOLATION: the only sanctioned
// producer of ErrKeybindingNotHandled from a chord handler is
// callKeybindingHandler in pkg/gui/keybindings.go:497-512, which
// short-circuits with NotHandled BEFORE invoking the underlying handler
// when GetDisabledReason returns AllowFurtherDispatching=true. That code
// path therefore never calls SetPendingChord.
//
// A handler that BOTH calls SetPendingChord AND returns NotHandled is
// outside the supported contract. The chord dispatcher does not defend
// against this case: pendingChord is cleared exactly once before the
// first handler invocation (gated by matchedExact), so a handler-set
// pendingChord persists into and past subsequent fall-through iterations.
//
// This test exists so that any future change to the dispatcher's
// pending-chord clearing strategy (e.g. re-clearing between iterations,
// or refusing to fall through after a SetPendingChord call) trips a
// failure here and forces a deliberate decision rather than a silent
// behavioral drift.
func TestChordDispatch_HandlerSetsPendingChordAndDeclines(t *testing.T) {
	g := &Gui{}
	view := NewView("commits", 0, 0, 10, 10, OutputNormal)
	g.views = append(g.views, view)

	secondFired := false
	g.keybindings = []*keybinding{
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(gg *Gui, _ *View) error {
			gg.SetPendingChord([]Key{NewKeyRune('x'), NewKeyRune('y')}, "commits")
			return ErrKeybindingNotHandled
		}),
		newChordKeybinding("commits", []Key{NewKeyRune('b'), NewKeyRune('b')}, func(*Gui, *View) error {
			secondFired = true
			return nil
		}),
	}

	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "commits"

	if err := g.execKeybindings(view, &GocuiEvent{Key: NewKeyRune('b')}); err != nil {
		t.Fatalf("execKeybindings returned error: %v", err)
	}

	// Fall-through still happens: the dispatcher treats NotHandled as
	// "decline, try next" regardless of side effects on pendingChord.
	if !secondFired {
		t.Fatal("expected second handler to fire after first declined via NotHandled")
	}

	// The pendingChord set by the first handler persists. matchedExact
	// gates the clear-on-first-match to a single ClearPendingChord call,
	// and that call happens BEFORE handler 1 runs. Handler 1's
	// SetPendingChord([x,y]) therefore lands after the clear, and no
	// subsequent iteration re-clears it.
	got := g.PendingChord()
	if len(got) != 2 || !got[0].Equals(NewKeyRune('x')) || !got[1].Equals(NewKeyRune('y')) {
		t.Fatalf("expected handler-set pendingChord [x, y] to persist past fall-through, got %v", got)
	}
	if g.PendingChordView() != "commits" {
		t.Fatalf("expected pendingChordView to be %q, got %q", "commits", g.PendingChordView())
	}
}
