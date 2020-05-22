# Concurrent Petri Net

Petri Net implementation based on golang concurrency patterns such as goroutines and channels

## Examples

Here're several example for Petri Net usage:
- [{ pin -> {t1 t2} -> pout }](./example/ptp/main.go) is an elementary network contains two places `in` and `out` and 
 concurrent transitions
- [{req -> echo -> res}}](./example/echo/README.md) is an implementation of HTTP echo server as a Petri Net 

## Benchmark

Solution overhead is about 3-5μs per transition

```
$: make bench
...
pkg: github.com/alxmsl/cpn/test
BenchmarkBlockPTP-4       	  455206	      2640 ns/op	     136 B/op	       5 allocs/op
BenchmarkBlockPTPTP-4     	  145244	      7762 ns/op	     304 B/op	       8 allocs/op
BenchmarkBlockPTPTPTP-4   	   83860	     15072 ns/op	     376 B/op	      10 allocs/op
BenchmarkQueuePTP-4       	  491024	      2337 ns/op	     136 B/op	       5 allocs/op
BenchmarkPPTP-4           	  384649	      3345 ns/op	     144 B/op	       5 allocs/op
BenchmarkPPTTP-4          	  358689	      3514 ns/op	     144 B/op	       5 allocs/op
```
