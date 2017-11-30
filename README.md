# Tenchmark

Thrift benchmark command line tools / framework.

**Contributions are welcomed**

## Usage
```
usage: tenchmark [<flags>] [<addr>]

Flags:
      --help                Show context-sensitive help (also try --help-long
                            and --help-man).
  -n, --requests=100        Number of requests to perform
  -c, --concurrency=2       Number of multiple requests to make at a time
      --path="/"            Http request path
      --protocol="binary"   Specify protocol factory
      --transport="socket"  Specify transport factory
      --transport-wrapper="buffered"
                            Specify transport wrapper
      --service=SERVICE     Specify service name

Args:
  [<addr>]  Server addr
```

## Results
```
Benchmarking :6000 (be patient)......

Server Address:         :6000

Concurrency level:      10
Time taken for tests:   0.010 seconds
Complete requests:      100
Failed requests:        0
Request per second:     10359.47 [#/sec] (mean)

Percentage of the requests served within a certain time (ms)
  50%     0.06
  66%     0.08
  75%     0.09
  80%     0.11
  90%     0.13
  95%     4.52
  98%     6.89
  99%     7.75
 100%     8.50 (longest request)
```
