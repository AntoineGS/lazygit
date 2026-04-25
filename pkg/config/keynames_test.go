package config

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/stretchr/testify/assert"
)

func TestKeyFromLabel_SingleKey(t *testing.T) {
	key, ok := KeyFromLabel("p")
	assert.True(t, ok)
	assert.True(t, key.IsSet())
	assert.Equal(t, "p", key.Str())
	assert.Empty(t, key.Rest())
}

func TestKeyFromLabel_Chord_TwoRunes(t *testing.T) {
	key, ok := KeyFromLabel("bp")
	assert.True(t, ok)
	assert.Equal(t, "b", key.Str())
	rest := key.Rest()
	assert.Len(t, rest, 1)
	assert.Equal(t, "p", rest[0].Str())
}

func TestKeyFromLabel_Chord_WithBracketedFirst(t *testing.T) {
	key, ok := KeyFromLabel("<c-g>p")
	assert.True(t, ok)
	assert.Equal(t, gocui.ModCtrl, key.Mod())
	assert.Equal(t, "g", key.Str())
	rest := key.Rest()
	assert.Len(t, rest, 1)
	assert.Equal(t, "p", rest[0].Str())
	assert.Equal(t, gocui.ModNone, rest[0].Mod())
}

func TestKeyFromLabel_Chord_AllBracketed(t *testing.T) {
	key, ok := KeyFromLabel("<c-g><c-p>")
	assert.True(t, ok)
	rest := key.Rest()
	assert.Len(t, rest, 1)
	assert.Equal(t, gocui.ModCtrl, rest[0].Mod())
	assert.Equal(t, "p", rest[0].Str())
}

func TestKeyFromLabel_Chord_Deep(t *testing.T) {
	key, ok := KeyFromLabel("abc")
	assert.True(t, ok)
	assert.Equal(t, "a", key.Str())
	rest := key.Rest()
	assert.Len(t, rest, 2)
	assert.Equal(t, "b", rest[0].Str())
	assert.Equal(t, "c", rest[1].Str())
}

func TestKeyFromLabel_Empty(t *testing.T) {
	key, ok := KeyFromLabel("")
	assert.True(t, ok)
	assert.False(t, key.IsSet())
}

func TestKeyFromLabel_Disabled(t *testing.T) {
	key, ok := KeyFromLabel("<disabled>")
	assert.True(t, ok)
	assert.False(t, key.IsSet())
}

func TestKeyFromLabel_RejectsUnterminatedBracket(t *testing.T) {
	_, ok := KeyFromLabel("b<bogus")
	assert.False(t, ok)
}

func TestKeyFromLabel_RejectsEscInNonFirstPosition(t *testing.T) {
	_, ok := KeyFromLabel("b<esc>c")
	assert.False(t, ok)
}

func TestKeyFromLabel_RejectsDisabledInsideChord(t *testing.T) {
	_, ok := KeyFromLabel("b<disabled>")
	assert.False(t, ok)
}
