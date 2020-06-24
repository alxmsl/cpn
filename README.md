# Concurrent Petri Net

Petri Net implementation based on golang concurrency patterns such as goroutines and channels

## Examples

Here're several example for Petri Net usage:
- [{req -> echo -> res}}](./example/echo/README.md) is an implementation of HTTP echo server as a Petri Net 
- [{ pin -> {t1 t2} -> pout }](./example/ptp/main.go) is an elementary network contains two places `in` and `out` and 
 concurrent transitions

## Benchmark

Solution overhead is about 3-5μs per transition

```
$: make bench
...
pkg: github.com/alxmsl/cpn/test
BenchmarkBlockPTP-4       	  374569	      2793 ns/op	     136 B/op	       5 allocs/op
BenchmarkBlockPTPTP-4     	  142035	      7753 ns/op	     304 B/op	       8 allocs/op
BenchmarkBlockPTPTPTP-4   	   90934	     14784 ns/op	     376 B/op	      10 allocs/op
BenchmarkQueuePTP-4       	  461224	      2444 ns/op	     136 B/op	       5 allocs/op
BenchmarkPTPP-4           	  113982	     10804 ns/op	     136 B/op	       5 allocs/op
BenchmarkPPTP-4           	  282396	      3828 ns/op	     144 B/op	       5 allocs/op
BenchmarkPPTTP-4          	  310167	      3467 ns/op	     144 B/op	       5 allocs/op
```
