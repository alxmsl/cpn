# HTTP echo server net

This is simple Petri Net which implements HTTP server

## Benchmark

```bash
wrk -c 5 -t 2 -d 10s --latency 'http://localhost:8080'
Running 10s test @ http://localhost:8080
  2 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   112.08us  300.56us  12.52ms   98.89%
    Req/Sec    20.29k     1.64k   22.30k    85.00%
  Latency Distribution
     50%   86.00us
     75%  101.00us
     90%  124.00us
     99%  436.00us
  403883 requests in 10.00s, 72.80MB read
Requests/sec:  40375.49
Transfer/sec:      7.28MB
```