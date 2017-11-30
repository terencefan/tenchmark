# Tenchmark

Thrift benchmark command line tools / framework.

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
