//go:build windows

package oscommands

import (
	"errors"
	"io"
)

// ansiStripReader wraps an io.Reader and removes ANSI escape sequences
// (CSI, OSC, and simple single-char escapes) from the byte stream so that
// lazygit's plain-text credential-prompt regexes match ConPTY output.
//
// The reader is streaming: a sequence split across Read() calls must not
// stall the consumer, because processOutput scans byte-by-byte and needs
// the terminating ':' of a credential prompt to appear immediately.
//
// At EOF, a lone ESC (state Esc) is emitted literally — it may be valid
// data from the source. Truncated CSI/OSC sequences are dropped rather
// than re-emitted, because flushing partial escape codes would put
// garbage into the command log. Worst case the detector simply does not
// match — same outcome as today's non-PTY path.
type ansiStripReader struct {
	src   io.Reader
	buf   [4096]byte
	state stripState
}

type stripState int

const (
	stateNormal stripState = iota
	stateEsc               // saw ESC, waiting to classify
	stateCSI               // inside ESC [ ... ; final byte is 0x40..0x7E
	stateOSC               // inside ESC ] ... ; terminated by BEL or ESC \
	stateOSCEsc            // inside OSC, saw ESC — next byte should be '\'
)

func newAnsiStripReader(src io.Reader) *ansiStripReader {
	return &ansiStripReader{src: src}
}

func (r *ansiStripReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Read at most len(p) bytes so we never grow the output beyond the caller's buffer.
	// Stripping is non-expanding: output size ≤ input size.
	n := min(len(p), len(r.buf))
	nr, err := r.src.Read(r.buf[:n])
	out := 0
	for i := range nr {
		b := r.buf[i]
		switch r.state {
		case stateNormal:
			if b == 0x1b { // ESC
				r.state = stateEsc
				continue
			}
			p[out] = b
			out++
		case stateEsc:
			switch b {
			case '[':
				r.state = stateCSI
			case ']':
				r.state = stateOSC
			default:
				// Simple escape (e.g. ESC 7 = save cursor). Swallow both bytes.
				r.state = stateNormal
			}
		case stateCSI:
			// Consume bytes until a final byte in 0x40..0x7E
			// (spec: parameter + intermediate bytes precede the final).
			if b >= 0x40 && b <= 0x7E {
				r.state = stateNormal
			}
			// else stay in CSI, consuming.
		case stateOSC:
			switch b {
			case 0x07: // BEL terminates
				r.state = stateNormal
			case 0x1b:
				r.state = stateOSCEsc
			}
			// else stay in OSC, consuming.
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
	// (state != stateNormal means we've dropped them). Return what we produced.
	// Malformed tail bytes are intentionally dropped per the comment at the top.
	// Exception: if the whole input was just an ESC, state==stateEsc at EOF —
	// we handle that by emitting it literally only if err == io.EOF.
	if errors.Is(err, io.EOF) && r.state == stateEsc {
		if out < len(p) {
			p[out] = 0x1b
			out++
			r.state = stateNormal
		}
	}
	return out, err
}
