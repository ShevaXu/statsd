package statsd

import (
	"math/rand"
	"net"
	"sync"
	"time"
)

// A Client is a StatsD client. It can be used by multiple goroutines.
type Client struct {
	buf    Flusher
	c      *config
	mu     sync.Mutex
	ticker *time.Ticker
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

// AutoFlush sets the client to flush the buffer every period.
// It override is set to true, this call will stop previous one
// and start a new one.
func (c *Client) AutoFlush(period time.Duration, override bool) {
	if period > 0 {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.ticker != nil {
			if !override {
				return
			}
			// stop the original ticker
			c.ticker.Stop()
		}
		// new ticker
		c.ticker = time.NewTicker(period)
		// do not block
		go func() {
			for _ = range c.ticker.C {
				c.buf.Flush()
			}
		}()
	}
}

// Packet size guildlines
// (from https://github.com/etsy/statsd/blob/master/docs/metric_types.md)
const (
	// These payload numbers take into account the maximum IP + UDP header sizes.
	PacketSizeFastEthernet      = 1432 // This is most likely for Intranets.
	PacketSizeGigabitEthernet   = 8932 // Jumbo frames can make use of this feature much more efficient.
	PacketSizeCommodityInternet = 512  // If you are routing over the internet ... (default)
)

// A config is the client's configurable options.
type config struct {
	addr          string
	maxPacketSize int
}

// An Option represents an option to change the default behaviours of
// the client, used as an argument to NewClient().
type Option func(*config)

// PacketSize sets the maximum bytes sent by the Client at a time.
func PacketSize(n int) Option {
	return func(c *config) {
		c.maxPacketSize = n
	}
}

// NewClient returns a new Client.
func NewClient(addr string, opts ...Option) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	c := &config{
		addr:          addr,
		maxPacketSize: PacketSizeCommodityInternet,
	}
	for _, opt := range opts {
		opt(c)
	}
	return &Client{
		buf: NewBuffer(c.maxPacketSize, 256, conn, IgnoreTrailingByte()),
	}, nil
}
