# Tenchmark

Thrift benchmark command line tools / framework.

**Contributions are welcomed**

We support following protocols and transports currently:

* protocol
    * binary
* transport
    * tcp socket
    * unix domain
* transport wrapper
    * framed
    * buffered

```
usage: main [<flags>] [<addr>]

Flags:
      --help                Show context-sensitive help (also try --help-long and --help-man).
  -n, --requests=10         Number of requests to perform
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

### Usage

* send ping request to :10010

```
tenchmark :10010
```

* send ping request to :10010 via framed transport

```
tenchmark :10010 --wrapper=framed
```

* send ping request via multiplexed protocol

```
tenchmark :10010 --service=<service_name>
```

* send ping request via unix domain socket

```
tenchmark /var/run/x.sock --transport=unix
```

### Advanced Usage

* send custom request

```
tenchmark -n 20000 -c 100 \
    --thrift-file ./example/test.thrift \
    --api-file ./example/test.json \
    --case case 1 \
    127.0.0.1:80
```

For further informations, see [examples]()

## Results
```
This is Tenchmark, Version 0.1
Copyright 2017 Terence Fan, Baixing, https://github.com/baixing
Licensed under the MIT

Benchmarking :6000 (be patient)......
Completed 1000 requests
Finished 1000 requests

Server Address:         :6000

Concurrency level:      10
Time taken for tests:   0.026 seconds
Complete requests:      1000
Failed requests:        0
Request per second:     37887.40 [#/sec] (mean)

Percentage of the requests served within a certain time (ms)
  50%     0.12
  66%     0.14
  75%     0.15
  80%     0.16
  90%     0.19
  95%     0.22
  98%     0.25
  99%     0.36
 100%    13.27 (longest request)
```
