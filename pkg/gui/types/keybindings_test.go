package types

import "testing"

func TestBinding_GetTooltip_FuncTakesPrecedence(t *testing.T) {
	b := &Binding{Tooltip: "static", TooltipFunc: func() string { return "dynamic" }}
	if got := b.GetTooltip(); got != "dynamic" {
		t.Fatalf("want dynamic, got %q", got)
	}
}

func TestBinding_GetTooltip_FallsBackToStatic(t *testing.T) {
	b := &Binding{Tooltip: "static"}
	if got := b.GetTooltip(); got != "static" {
		t.Fatalf("want static, got %q", got)
	}
}

func TestBinding_GetTooltip_EmptyWhenNothingSet(t *testing.T) {
	b := &Binding{}
	if got := b.GetTooltip(); got != "" {
		t.Fatalf("want empty, got %q", got)
	}
}

func TestBinding_HiddenInChordPopup_NilDefaultsFalse(t *testing.T) {
	b := &Binding{}
	if b.IsHiddenInChordPopup() {
		t.Fatal("nil HiddenInChordPopup should not hide")
	}
}

func TestBinding_HiddenInChordPopup_FuncDecides(t *testing.T) {
	b := &Binding{HiddenInChordPopup: func() bool { return true }}
	if !b.IsHiddenInChordPopup() {
		t.Fatal("HiddenInChordPopup returning true should hide")
	}
}
