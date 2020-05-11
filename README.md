# Concurrent Petri Net

Petri Net implementation based on golang concurrency patterns such as goroutines and channels

## Examples

Here're several example for Petri Net usage:
- [Simple Petri Net](./example/ptp/main.go) is an elementary network contains two places `in` and `out` and just one 
 transition
- [Echo server](./example/echo/README.md) is an implementation of HTTP echo server as a Petri Net 

## Benchmark

Solution overhead is about 3-5μs per transition

```
$: make bench
...
pkg: github.com/alxmsl/cpn/test
BenchmarkBlockPTP-4       	  387151	      2863 ns/op	     136 B/op	       5 allocs/op
BenchmarkBlockPTPTP-4     	  163048	      7701 ns/op	     304 B/op	       8 allocs/op
BenchmarkBlockPTPTPTP-4   	   88299	     13658 ns/op	     376 B/op	      10 allocs/op
```
