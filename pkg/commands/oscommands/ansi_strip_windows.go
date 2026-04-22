//go:build windows

package oscommands

import (
	"errors"
	"io"
)

// Allows removing ANSI escape sequences that are from a
// multi-Read() call, where a regular regex implementation would fail to
// match the full sequence.

type ansiStripReader struct {
	src   io.Reader
	buf   [4096]byte
	state stripState
}

type stripState int

const (
	stateNormal stripState = iota
	stateEsc               // saw ESC, waiting to classify
	stateCSI               // inside ESC [ ... ; final byte is csiFinalMin..csiFinalMax
	stateOSC               // inside ESC ] ... ; terminated by BEL or ESC \
	stateOSCEsc            // inside OSC, saw ESC — next byte should be '\'
)

const (
	esc         byte = 0x1B // ESC
	bel         byte = 0x07 // BEL — OSC string terminator (alternative to ESC \)
	csiFinalMin byte = 0x40 // '@' — lowest CSI final byte
	csiFinalMax byte = 0x7E // '~' — highest CSI final byte
)

func newAnsiStripReader(src io.Reader) *ansiStripReader {
	return &ansiStripReader{src: src}
}

func (r *ansiStripReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	nDst := min(len(p), len(r.buf))
	nSrc, err := r.src.Read(r.buf[:nDst])
	n := 0
	for i := range nSrc {
		b := r.buf[i]
		switch r.state {
		case stateNormal:
			if b == esc {
				r.state = stateEsc
				continue
			}
			p[n] = b
			n++
		case stateEsc:
			switch b {
			case '[':
				r.state = stateCSI
			case ']':
				r.state = stateOSC
			default:
				r.state = stateNormal
			}
		case stateCSI:
			// Consume bytes until a final byte in csiFinalMin..csiFinalMax
			if b >= csiFinalMin && b <= csiFinalMax {
				r.state = stateNormal
			}
		case stateOSC:
			switch b {
			case bel:
				r.state = stateNormal
			case esc:
				r.state = stateOSCEsc
			}
		case stateOSCEsc:
			// Expect '\' to complete ST. Anything else: treat as re-entry to OSC.
			if b == '\\' {
				r.state = stateNormal
			} else {
				r.state = stateOSC
			}
		}
	}
	// If upstream EOF arrived mid-sequence, we have already consumed those bytes
	// (state != stateNormal means we've dropped them).
	// Malformed tail bytes are intentionally dropped.
	if errors.Is(err, io.EOF) && r.state == stateEsc {
		if n < len(p) {
			p[n] = esc
			n++
			r.state = stateNormal
		}
	}
	return n, err
}
