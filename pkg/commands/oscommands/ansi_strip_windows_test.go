//go:build windows

package oscommands

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestAnsiStripReader(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain ascii passes through",
			input: "hello world\n",
			want:  "hello world\n",
		},
		{
			name:  "strip simple CSI SGR",
			input: "\x1b[32mgreen\x1b[0m",
			want:  "green",
		},
		{
			name:  "strip CSI clear screen",
			input: "before\x1b[2Jafter",
			want:  "beforeafter",
		},
		{
			name:  "strip CSI with multiple parameters",
			input: "\x1b[1;31;4mbold-red-underline\x1b[0m",
			want:  "bold-red-underline",
		},
		{
			name:  "strip CSI private mode",
			input: "\x1b[?25lhidden\x1b[?25h",
			want:  "hidden",
		},
		{
			name:  "strip OSC terminated by BEL",
			input: "before\x1b]0;window title\x07after",
			want:  "beforeafter",
		},
		{
			name:  "strip OSC terminated by ST",
			input: "before\x1b]0;window title\x1b\\after",
			want:  "beforeafter",
		},
		{
			name:  "strip simple escape (cursor save)",
			input: "before\x1b7after",
			want:  "beforeafter",
		},
		{
			name:  "credential prompt survives color wrap",
			input: "\x1b[32mEnter passphrase for key '/c/Users/x/.ssh/id_ed25519':\x1b[0m",
			want:  "Enter passphrase for key '/c/Users/x/.ssh/id_ed25519':",
		},
		{
			name:  "credential prompt with extra CSI noise",
			input: "\x1b[?25h\x1b[0mEnter passphrase for key 'id_rsa': \x1b[?25l",
			want:  "Enter passphrase for key 'id_rsa': ",
		},
		{
			name:  "lone ESC at EOF emitted literally",
			input: "hello\x1b",
			want:  "hello\x1b",
		},
		{
			name:  "malformed CSI (EOF mid-sequence) does not swallow prior bytes",
			input: "hello\x1b[31",
			want:  "hello",
		},
		{
			name:  "ESC followed by ordinary char is treated as escape",
			input: "a\x1bDb",
			want:  "ab",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newAnsiStripReader(strings.NewReader(tc.input))
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("ReadAll: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %q want %q", string(got), tc.want)
			}
		})
	}
}

func TestAnsiStripReader_ByteByByte(t *testing.T) {
	// Feed the reader one byte at a time — sequences must not get stuck
	// across Read boundaries.
	in := "\x1b[32mEnter passphrase for key 'x': \x1b[0m"
	want := "Enter passphrase for key 'x': "

	r := newAnsiStripReader(&oneByteAtATime{s: in})
	var out bytes.Buffer
	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			out.Write(buf[:n])
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
	}
	if out.String() != want {
		t.Errorf("got %q want %q", out.String(), want)
	}
}

func TestAnsiStripReader_ByteByByte_OSC(t *testing.T) {
	// OSC with ST terminator split across Read boundaries.
	in := "pre\x1b]0;title\x1b\\post"
	want := "prepost"

	r := newAnsiStripReader(&oneByteAtATime{s: in})
	var out bytes.Buffer
	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			out.Write(buf[:n])
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
	}
	if out.String() != want {
		t.Errorf("got %q want %q", out.String(), want)
	}
}

// oneByteAtATime returns bytes one at a time from a string, then io.EOF.
type oneByteAtATime struct {
	s string
	i int
}

func (r *oneByteAtATime) Read(p []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = r.s[r.i]
	r.i++
	return 1, nil
}
