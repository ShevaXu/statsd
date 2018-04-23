package statsd_test

import (
	"bytes"
	"testing"

	"github.com/ShevaXu/golang/assert"
	"github.com/ShevaXu/statsd"
)

var (
	hello = []byte("hello")
	world = []byte("world")
)

func TestBuffer(t *testing.T) {
	a := assert.NewAssert(t)
	var out bytes.Buffer
	buf := statsd.NewBuffer(10, 5, &out)

	buf.Start()
	buf.Write(hello)
	buf.End()

	buf.Flush()
	a.Equal(hello, out.Bytes(), "force flush")
	out.Reset()

	buf.Start()
	buf.Write(hello)
	buf.End()

	buf.Start()
	buf.AppendByte(' ')
	buf.Write(world)
	buf.End()
	a.Equal(hello, out.Bytes(), "auto flush")
	out.Reset()

	buf.Flush()
	a.Equal(" world", out.String(), "flush the rest")
	out.Reset()

	buf.Flush()
	a.Equal(0, len(out.Bytes()), "flush when empty")
}

func TestBufferIgnoreTrailingByte(t *testing.T) {
	a := assert.NewAssert(t)
	var out bytes.Buffer
	buf := statsd.NewBuffer(10, 5, &out, statsd.IgnoreTrailingByte())

	buf.Start()
	buf.Write(hello)
	buf.AppendByte(' ')
	buf.End()

	buf.Start()
	buf.Write(world)
	buf.End()
	a.Equal(hello, out.Bytes(), "flush without trailing byte")
}
