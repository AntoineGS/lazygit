package gocui

import "testing"

func TestSuppressChordClear_BlocksClearWhenSet(t *testing.T) {
	g := &Gui{}
	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "files"

	g.SuppressChordClear(true)
	g.ClearPendingChord()

	if g.pendingChord == nil {
		t.Fatal("ClearPendingChord must be a no-op while suppressChordClear is true")
	}

	g.SuppressChordClear(false)
	g.ClearPendingChord()
	if g.pendingChord != nil {
		t.Fatal("ClearPendingChord should clear once suppress is off")
	}
}

func TestSetPendingChordView_UpdatesViewName(t *testing.T) {
	g := &Gui{}
	g.pendingChord = []Key{NewKeyRune('b')}
	g.pendingChordView = "files"

	g.SetPendingChordView("localBranches")

	if g.pendingChordView != "localBranches" {
		t.Fatalf("expected pendingChordView=localBranches, got %q", g.pendingChordView)
	}
}

func TestSetGlobalChordPrefixes_StoresAndDefensiveCopies(t *testing.T) {
	g := &Gui{}
	src := []Key{NewKeyRune('f'), NewKeyRune('r')}
	g.SetGlobalChordPrefixes(src)

	if len(g.globalChordPrefixes) != 2 {
		t.Fatalf("expected 2 stored prefixes, got %d", len(g.globalChordPrefixes))
	}

	src[0] = NewKeyRune('z')
	if !g.globalChordPrefixes[0].Equals(NewKeyRune('f')) {
		t.Fatal("SetGlobalChordPrefixes must defensively copy its input")
	}

	g.SetGlobalChordPrefixes(nil)
	if len(g.globalChordPrefixes) != 0 {
		t.Fatalf("expected empty after nil set, got %d", len(g.globalChordPrefixes))
	}
}
