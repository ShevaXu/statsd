package statsd

import (
	"io"
	"strconv"
	"sync"
)

// A Buffer is a (supposedly) fix-sized Buffer of bytes with auto
// out-of-capacity flushing (to a io.Writer), mainly used for sending
// multiple messages to the same destination. The zero value for
// Buffer will panic when flush due to nil io.Writer. Use NewBuffer()
// to init instead.
type Buffer struct {
	// options
	cap                int
	ignoreTrailingByte bool
	// status
	safeLen int
	// protected
	mu sync.Mutex
	w  io.Writer
	v  []byte
}

// Write appends the bytes to the buffer.
func (b *Buffer) Write(s []byte) {
	b.v = append(b.v, s...)
}

// WriteString appends a string to the buffer.
func (b *Buffer) WriteString(s string) {
	b.v = append(b.v, s...)
}

// AppendByte appends a single byte to the buffer.
func (b *Buffer) AppendByte(c byte) {
	b.v = append(b.v, c)
}

// AppendInt converts a int to bytes and appends to the buffer
// without extra allocation.
func (b *Buffer) AppendInt(i int64) {
	b.v = strconv.AppendInt(b.v, i, 10)
}

// AppendFloat converts a float to bytes and appends to the buffer
// without extra allocation.
func (b *Buffer) AppendFloat(f float64) {
	b.v = strconv.AppendFloat(b.v, f, 'f', -1, 64)
}

// flush writes the first l bytes to the underlying io.Writer, and
// shifts the remaining bytes up. It flushes all if l = 0.
func (b *Buffer) flush(l int) {
	n := len(b.v)
	if n == 0 {
		return
	}
	if l == 0 {
		// flush all
		l = n
	}
	// ignore write error
	if b.ignoreTrailingByte {
		// trim '\n'
		b.w.Write(b.v[:l-1])
	} else {
		b.w.Write(b.v[:l])
	}
	// reset
	if l < n {
		copy(b.v, b.v[l:])
	}
	b.v = b.v[:n-l]
}

// Start locks the buffer to enable a series of transactional writes;
// it must be followed by a End() call when finish writing. This allows
// concurrent access to the buffer by different goroutines. It also
// marks last safe len below capacity for auto-flushing.
func (b *Buffer) Start() {
	b.mu.Lock()
	b.safeLen = len(b.v)
}

// End finishes the transactional writes by unlocking the buffer. It
// triggers flush if the writes cause buffer' length > capacity.
func (b *Buffer) End() {
	if len(b.v) >= b.cap {
		b.flush(b.safeLen)
	}
	b.mu.Unlock()
}

// Flush forces the buffer to write all its bytes immediately.
func (b *Buffer) Flush() {
	b.mu.Lock()
	b.flush(0)
	b.mu.Unlock()
}

// WriteOnce writes to the buffer as a complete transaction. It
// is the same as call just one Write() between Start() & End().
func (b *Buffer) WriteOnce(s []byte) {
	b.Start()
	b.v = append(b.v, s...)
	b.End()
}

// WriteStringOnce writes a string as a complete transaction.
func (b *Buffer) WriteStringOnce(s string) {
	b.Start()
	b.v = append(b.v, s...)
	b.End()
}

// A BufferOption represents the Buffer's functional option, used
// as an argument to NewBuffer().
type BufferOption func(*Buffer)

// IgnoreTrailingByte returns a BufferOptions for setting the Buffer
// to ignore the trailing byte when flushing; it is useful for trimming
// the last separator (e.g., '\n') for multiple messages.
func IgnoreTrailingByte() BufferOption {
	return func(b *Buffer) {
		b.ignoreTrailingByte = true
	}
}

// NewBuffer returns a Buffer. The underlying bytes array is supposed to
// be initialized once and never grows again, if the cap & cache size is
// set right. The cap is related to the output/writing pattern, e.g., max
// network packet size to avoid framentation; and cache should fits a single
// message.
func NewBuffer(cap, cache int, w io.Writer, opts ...BufferOption) *Buffer {
	b := &Buffer{
		cap: cap,
		w:   w,
		v:   make([]byte, 0, cap+cache),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}
