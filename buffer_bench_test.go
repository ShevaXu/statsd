package statsd_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"github.com/ShevaXu/statsd"
)

var (
	testMetricStr = "benchmark.statsd:100|c@0.1"
	testMetricBs  = []byte(testMetricStr)
	udpPacketSize = 1440
)

func BenchmarkBufferPacketSize(b *testing.B) {
	buf := statsd.NewBuffer(udpPacketSize, 100, ioutil.Discard)

	for i := 0; i < b.N; i++ {
		buf.WriteOnce(testMetricBs)
	}
}

func BenchmarkBufferCast(b *testing.B) {
	buf := statsd.NewBuffer(udpPacketSize, 100, ioutil.Discard)

	for i := 0; i < b.N; i++ {
		buf.WriteOnce([]byte(testMetricStr))
	}
}

func BenchmarkStdBuffer(b *testing.B) {
	var buf bytes.Buffer
	var mu sync.Mutex

	for i := 0; i < b.N; i++ {
		mu.Lock()
		if buf.Len()+len(testMetricBs) > udpPacketSize {
			ioutil.Discard.Write(buf.Bytes())
			buf.Reset()
		}
		buf.Write(testMetricBs)
		mu.Unlock()
	}
}

func BenchmarkStdBufferToString(b *testing.B) {
	var buf bytes.Buffer
	var mu sync.Mutex

	for i := 0; i < b.N; i++ {
		mu.Lock()
		if buf.Len()+len(testMetricStr) > udpPacketSize {
			_ = buf.String() // this method has allocations
			buf.Reset()
		}
		buf.WriteString(testMetricStr)
		mu.Unlock()
	}
}

func BenchmarkStdBuilder(b *testing.B) {
	var buf strings.Builder
	var mu sync.Mutex

	for i := 0; i < b.N; i++ {
		mu.Lock()
		if buf.Len()+len(testMetricStr) > udpPacketSize {
			_ = buf.String() // this method has allocations
			buf.Reset()
		}
		buf.WriteString(testMetricStr)
		mu.Unlock()
	}
}
