package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// chordPreview renders cmd through the given style. It returns "" when
// no git command is available so unbound bindings don't show stale
// previews.
func chordPreview(c types.IGuiCommon, s style.TextStyle, cmd func() string) string {
	if c.Git() == nil {
		return ""
	}
	return s.Sprint(cmd())
}
