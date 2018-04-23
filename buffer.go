package statsd

import (
	"io"
	"sync"
)

// Buffer is a re-usable buffer for
// writing to a destination (io.Writer).
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

func (b *Buffer) Write(s []byte) {
	b.v = append(b.v, s...)
}

func (b *Buffer) AppendByte(c byte) {
	b.v = append(b.v, c)
}

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

func (b *Buffer) Start() {
	b.mu.Lock()
	b.safeLen = len(b.v)
}

func (b *Buffer) End() {
	if len(b.v) >= b.cap {
		b.flush(b.safeLen)
	}
	b.mu.Unlock()
}

func (b *Buffer) Flush() {
	b.mu.Lock()
	b.flush(0)
	b.mu.Unlock()
}

type BufferOption func(*Buffer)

func IgnoreTrailingByte() BufferOption {
	return func(b *Buffer) {
		b.ignoreTrailingByte = true
	}
}

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
