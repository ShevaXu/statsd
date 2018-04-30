package statsd

import (
	"math/rand"
	"net"
)

// A Client is a StatsD client. It can be used by multiple goroutines.
type Client struct {
	buf Flusher
}

// metric appends bytes to the buffer according to the StatsD protocal.
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

// Gauge records a int value calculated at client-side.
func (c *Client) Gauge(bucket string, v int64) {
	// bucket:123|g
	c.metric(bucket, v, 'g')
}

// Counting adds a int value to the bucket (to be summerized at server-side).
func (c *Client) Counting(bucket string, v int64) {
	// bucket:1|c
	c.metric(bucket, v, 'c')
}

// Timing records a event timing in millisecond.
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

// Sampling increments the counter every "rate" (usually <1) of the time.
// It allows hot-paths to send less often.
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
