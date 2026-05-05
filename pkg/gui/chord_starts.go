package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// New chords are allowed only when no popup is open. With the chord
// menu open, dispatch goes through the menu's per-item keybindings;
// starting a fresh gocui chord here would race with that and reset
// the prefix. The chord menu uses a TEMPORARY_POPUP context, so it
// already appears in popupKeys when open — no separate flag needed.
func chordStartsEnabled(popupKeys []types.ContextKey) bool {
	return len(popupKeys) == 0
}
