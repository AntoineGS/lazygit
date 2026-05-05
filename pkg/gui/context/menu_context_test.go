package context

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/theme"
)

func TestRenderDisabledColumns_StylesEveryColumn(t *testing.T) {
	cols := []string{"Label", "Tooltip preview"}
	got := renderDisabledColumns(cols)

	if len(got) != 2 {
		t.Fatalf("want 2 columns, got %d", len(got))
	}

	// Probe the styled prefix produced by DisabledTextStyle+strikethrough so
	// we can detect it on every column. If colour is disabled in the test env
	// (no ANSI escape produced), there's nothing to assert; skip.
	probe := theme.DisabledTextStyle.SetStrikethrough().Sprint("X")
	idx := strings.Index(probe, "X")
	if idx <= 0 {
		t.Skip("colour disabled in test env")
	}
	styledOpen := probe[:idx]

	for i, col := range got {
		if !strings.HasPrefix(col, styledOpen) {
			t.Fatalf("column %d not greyed+struck: %q", i, col)
		}
	}
}
