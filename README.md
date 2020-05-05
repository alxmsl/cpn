# Real-time Petri Net

Real-time Petri Net implementation based on set of goroutines and channels

## Benchmark

Solution overhead is about 3-5μs per transition

```
$: make bench
...
pkg: github.com/alxmsl/rtpn/test
BenchmarkBlockPTP-4       	  387151	      2863 ns/op	     136 B/op	       5 allocs/op
BenchmarkBlockPTPTP-4     	  163048	      7701 ns/op	     304 B/op	       8 allocs/op
BenchmarkBlockPTPTPTP-4   	   88299	     13658 ns/op	     376 B/op	      10 allocs/op
```
