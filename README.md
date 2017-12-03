# Tenchmark

Thrift benchmark command line tool.

**Contributions are welcomed**

We support following protocols and transports currently:

* protocol
    * binary
* transport
    * tcp socket
    * unix domain
    * http
* transport wrapper
    * framed
    * buffered

```
Thrift benchmark command-line tools

Flags:
  --help               Show context-sensitive help (also try --help-long and
                       --help-man).
  --protocol="binary"  Specify protocol factory
  --service=SERVICE    Specify service name (multiplexed)

Commands:
  help [<command>...]
    Show help.

  run [<flags>] [<addr>]
    Run benchmark tests

  build --json=JSON [<flags>] <thrift>
    Build cases from thrift file and json inputs
```

### Usage

* send ping request to :10010

```
$ tenchmark run :10010
```

* send ping request with multiplexed protocol

```
$ tenchmark run :10010 --service=<service_name>
```

* send ping request to :10010 via framed transport

```
$ tenchmark run :10010 --wrapper=framed
```

* send ping request via unix domain socket

```
$ tenchmark run /var/run/x.sock --transport=unix
```

* send ping request via http

```
$ tenchmark run http://<host>:<port>/<path> --transport=http
```

### Advanced Usage

You can build your own cases through **tenchmark build** command.

```
$ tenchmark build example/ping.thrift --json=example/api.json
cases/case1.in sucessfully generated.
```

And you can specify the case you've built through following command.

```
$ tenchmark run --case=cases/case1.in
```

For further informations, see [examples (work in progress)]()

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
