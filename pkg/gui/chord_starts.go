package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// New chords are allowed only when no popup is open. With the chord
// menu open, dispatch goes through the menu's per-item keybindings;
// starting a fresh gocui chord here would race with that and reset
// the prefix.
func chordStartsEnabled(popupKeys []types.ContextKey, _ bool) bool {
	return len(popupKeys) == 0
}
