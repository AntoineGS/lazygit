package oscommands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTailBuffer(t *testing.T) {
	tests := []struct {
		name   string
		cap    int
		writes []string
		want   string
	}{
		{"empty", 10, nil, ""},
		{"single fits", 10, []string{"hello"}, "hello"},
		{"single exactly fills", 5, []string{"hello"}, "hello"},
		{"single overflows", 3, []string{"hello"}, "llo"},
		{"two writes fit", 10, []string{"hi", "lo"}, "hilo"},
		{"slide on second", 4, []string{"abc", "de"}, "bcde"},
		{"big then small", 4, []string{"hello", "x"}, "llox"},
		{"many small slides", 3, []string{"a", "b", "c", "d", "e"}, "cde"},
		{"empty write", 5, []string{""}, ""},
		{"single far exceeds cap", 2, []string{"abcdef"}, "ef"},
		{"second write exactly cap", 3, []string{"xy", "abc"}, "abc"},
		{"second write far exceeds cap", 3, []string{"xy", "abcdef"}, "def"},
		{"two writes exactly fill", 5, []string{"ab", "cde"}, "abcde"},
		{"big write into full buffer", 4, []string{"abcd", "wxyz"}, "wxyz"},
		{"oversize write into full buffer", 4, []string{"abcd", "vwxyz"}, "wxyz"},
		{"cap one single char", 1, []string{"a"}, "a"},
		{"cap one overflows", 1, []string{"abc"}, "c"},
		{"cap one many writes", 1, []string{"a", "b", "c"}, "c"},
		{"empty then write", 4, []string{"", "abc"}, "abc"},
		{"write between non-empty", 4, []string{"ab", "", "cd"}, "abcd"},
		{"wrap-around overshoot", 5, []string{"abcde", "fghi", "jk"}, "ghijk"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTailBuffer(tt.cap)
			for _, w := range tt.writes {
				n, err := b.Write([]byte(w))
				require.NoError(t, err)
				require.Equal(t, len(w), n)
			}
			assert.Equal(t, tt.want, b.String())
			assert.Equal(t, tt.want, string(b.Bytes()))
		})
	}
}

func TestTailBufferStaysBounded(t *testing.T) {
	const cap = 1024
	b := newTailBuffer(cap)
	chunk := strings.Repeat("x", 4096)
	for i := 0; i < 10000; i++ {
		_, _ = b.Write([]byte(chunk))
		require.LessOrEqual(t, len(b.Bytes()), cap)
	}
	assert.Equal(t, strings.Repeat("x", cap), b.String())
}
