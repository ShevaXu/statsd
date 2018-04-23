package statsd

import (
	"io"
)

// Buffer is a re-usable buffer for
// writing to a destination (io.Writer).
type Buffer struct {
	w                  io.Writer
	v                  []byte
	cap                int
	ignoreTrailingByte bool
}

func (b *Buffer) Write(s []byte) {
	b.v = append(b.v, s...)
}

func (b *Buffer) WriteByte(c byte) {
	b.v = append(b.v, c)
}

func (b *Buffer) IsFull() bool {
	return len(b.v) >= b.cap
}

func (b *Buffer) IsEmpty() bool {
	return len(b.v) == 0
}

func (b *Buffer) Flush(l int) {
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

type BufferOption func(*Buffer)

func IgnoreTrailingByte() BufferOption {
	return func(b *Buffer) {
		b.ignoreTrailingByte = true
	}
}

func NewBuffer(cap, cache int, w io.Writer, opts ...BufferOption) *Buffer {
	b := &Buffer{
		w: w,
		v: make([]byte, 0, cap+cache),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}
