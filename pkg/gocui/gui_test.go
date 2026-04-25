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
