package statsd_test

import (
	"bytes"
	"testing"

	"github.com/ShevaXu/golang/assert"
	"github.com/ShevaXu/statsd"
)

func TestBuffer(t *testing.T) {
	a := assert.NewAssert(t)
	var out bytes.Buffer
	buf := statsd.NewBuffer(10, 5, &out)
	str := []byte("hello world")

	buf.Write(str)
	a.True(buf.IsFull(), "is full")
	buf.Flush(0)
	a.Equal(str, out.Bytes(), "flush all")
	a.True(buf.IsEmpty(), "is empty")

	out.Reset()
	buf.Flush(0)
	a.Equal(0, len(out.Bytes()), "flush nothing if empty")

	buf.Write(str)
	buf.Flush(6)
	a.Equal(str[:6], out.Bytes(), "partial flush")

	out.Reset()
	buf.Flush(0)
	a.Equal(str[6:], out.Bytes(), "flush the rest")

	for _, c := range str {
		buf.WriteByte(c)
	}
	out.Reset()
	buf.Flush(0)
	a.Equal(str, out.Bytes(), "flush all byte by byte")
}

func TestBufferIgnoreTrailingByte(t *testing.T) {
	a := assert.NewAssert(t)
	var out bytes.Buffer
	buf := statsd.NewBuffer(10, 5, &out, statsd.IgnoreTrailingByte())
	str := []byte("hello world\n")

	buf.Write(str)
	buf.Flush(0)
	a.Equal(str[:len(str)-1], out.Bytes(), "flush all but last byte")
}
