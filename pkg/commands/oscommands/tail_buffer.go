package oscommands

const cmdOutputTailCap = 10240 * 1024

// io.Writer that retains only the last `cap` bytes written, to control
// memory usage when capturing command output, uses a ring buffer approach.
// Not safe for concurrent writes (matches bytes.Buffer's semantics).
type tailBuffer struct {
	buf   []byte
	cap   int
	start int
}

func newTailBuffer(cap int) *tailBuffer {
	return &tailBuffer{cap: cap, start: 0}
}

// The returned length can be > cap as some of the data may have been trimmed
func (b *tailBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if n == 0 {
		return 0, nil
	}
	// New data alone fills the buffer; only its tail matters.
	if n >= b.cap {
		b.buf = append(b.buf[:0], p[n-b.cap:]...)
		b.start = 0
		return n, nil
	}
	// Still filling for the first time: grow b.buf via append. b.start is 0
	// here, since it can only advance after the buffer has filled to cap.
	if len(b.buf) < b.cap {
		room := b.cap - len(b.buf)
		if n <= room {
			b.buf = append(b.buf, p...)
			return n, nil
		}
		b.buf = append(b.buf, p[:room]...)
		p = p[room:]
	}
	// Ring mode: buffer is full, overwrite at b.start, wrapping past the end.
	if b.start+len(p) <= b.cap {
		copy(b.buf[b.start:], p)
		b.start += len(p)
		if b.start == b.cap {
			b.start = 0
		}
	} else {
		first := b.cap - b.start
		copy(b.buf[b.start:], p[:first])
		copy(b.buf, p[first:])
		b.start = len(p) - first
	}
	return n, nil
}

func (b *tailBuffer) Bytes() []byte {
	if b.start == 0 {
		return b.buf
	}
	return append(b.buf[b.start:], b.buf[:b.start]...)
}

func (b *tailBuffer) String() string {
	return string(b.Bytes())
}
