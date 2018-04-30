# statsd

> StatsD is a simple daemon for aggregating application metrics.

*StatsD* was initialized by Flickr and *Cal Henderson* posted about it, [*Counting and timing*](http://code.flickr.com/blog/2008/10/27/counting-timing/), with the Perl [code](https://github.com/iamcal/Flickr-StatsD). 
Now [etsy/statsd](https://github.com/etsy/statsd) (Node.js) is the main reference repo.

## Benchmarks

### Buffer

```sh
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/ShevaXu/statsd
BenchmarkBufferPacketSize-8    	50000000	        25.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBufferCast-8          	50000000	        31.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdBuffer-8           	100000000	        22.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdBufferToString-8   	50000000	        28.6 ns/op	      27 B/op	       0 allocs/op
BenchmarkStdBuilder-8          	30000000	        40.3 ns/op	      73 B/op	       0 allocs/op
```

### Client

```sh
# cd benchmark
$ go test -bench . -benchmem
BenchmarkStatsd-8       	 3000000	       450 ns/op	       0 B/op	       0 allocs/op
BenchmarkAlexcesaro-8   	 3000000	       432 ns/op	       0 B/op	       0 allocs/op
```

## Refs

- blog [*Measure Anything Measure Everything*](https://codeascraft.com/2011/02/15/measure-anything-measure-everything/)
- go-client [alexcesaro/statsd](https://github.com/alexcesaro/statsd)
