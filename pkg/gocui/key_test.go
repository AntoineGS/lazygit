// Copyright 2026 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "testing"

func TestSequence_StripsRest(t *testing.T) {
	g := NewKeyRune('g')
	p := NewKeyRune('p')
	x := NewKeyRune('x')
	chord := g.WithRest([]Key{p, x})

	if !chord.HasRest() {
		t.Fatal("setup error: chord should HasRest before Sequence()")
	}

	seq := chord.Sequence()
	if len(seq) != 3 {
		t.Fatalf("expected sequence length 3, got %d", len(seq))
	}
	for i, k := range seq {
		if k.HasRest() {
			t.Errorf("Sequence()[%d] still has Rest set; the documented contract is that Sequence strips it", i)
		}
	}

	if !seq[0].Equals(g) {
		t.Errorf("Sequence()[0] should equal head key 'g'; got %+v", seq[0])
	}
}

func TestSequence_SingleKey_OneElement(t *testing.T) {
	k := NewKeyRune('q')
	seq := k.Sequence()
	if len(seq) != 1 {
		t.Fatalf("expected length 1 for single-key, got %d", len(seq))
	}
	if seq[0].HasRest() {
		t.Error("single-key Sequence()[0] should have no rest")
	}
}
