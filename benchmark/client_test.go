package statsdbench

import (
	"net"
	"testing"
	"time"

	"github.com/ShevaXu/statsd"
	ac "gopkg.in/alexcesaro/statsd.v2"
)

const (
	addr        = ":0"
	prefix      = "prefix."
	prefixNoDot = "prefix"
	counterKey  = "foo.bar.counter"
	gaugeKey    = "foo.bar.gauge"
	gaugeValue  = 42
	timingKey   = "foo.bar.timing"
	timingValue = 153 * time.Millisecond
	flushPeriod = 100 * time.Millisecond
)

func BenchmarkStatsd(b *testing.B) {
	s := newServer()
	c, err := statsd.NewClient(s.LocalAddr().String())
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Counting(counterKey, 1)
		c.Gauge(gaugeKey, gaugeValue)
		c.Timing(timingKey, int64(timingValue))
	}
	// c.Close()
	s.Close()
}

func BenchmarkAlexcesaro(b *testing.B) {
	s := newServer()
	c, err := ac.New(
		ac.Address(s.LocalAddr().String()),
		ac.Prefix(prefixNoDot),
		ac.FlushPeriod(flushPeriod),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Increment(counterKey)
		c.Gauge(gaugeKey, gaugeValue)
		c.Timing(timingKey, timingValue)
	}
	c.Close()
	s.Close()
}

func newServer() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		buf := make([]byte, 512)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				return
			}
		}
	}()
	return conn
}
