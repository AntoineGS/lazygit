package gocui

import "testing"

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
