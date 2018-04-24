package statsd

import (
	"math/rand"
	"net"
)

type Client struct {
	buf Flusher
}

func (c *Client) metric(bucket string, v int64, typ byte) {
	c.buf.Start()
	c.buf.WriteString(bucket)
	c.buf.AppendByte(':')
	c.buf.AppendInt(v)
	c.buf.AppendByte('|')
	c.buf.AppendByte(typ)
	c.buf.AppendByte('\n')
	c.buf.End()
}

func (c *Client) Gauge(bucket string, v int64) {
	// bucket:123|g
	c.metric(bucket, v, 'g')
}

func (c *Client) Counting(bucket string, v int64) {
	// bucket:1|c
	c.metric(bucket, v, 'c')
}

func (c *Client) Timing(bucket string, v int64) {
	// bucket:3600|ms
	c.buf.Start()
	c.buf.WriteString(bucket)
	c.buf.AppendByte(':')
	c.buf.AppendInt(v)
	c.buf.WriteString("|ms")
	c.buf.AppendByte('\n')
	c.buf.End()
}

func (c *Client) Sampling(bucket string, rate float64) {
	// bucket:1|c|@0.1
	if rand.Float64() > rate {
		return
	}
	c.buf.Start()
	c.buf.WriteString(bucket)
	c.buf.WriteString("1|c@")
	c.buf.AppendFloat(rate)
	c.buf.AppendByte('\n')
	c.buf.End()
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		buf: NewBuffer(1440, 200, conn, IgnoreTrailingByte()),
	}, nil
}
