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
-f, --file                Path to thrift file
--api                 Path to api file
-c  --case                Case

Args:
[<addr>]  Server addr
```

### Usage
```
tenchmark -n 20000 -c 100 --transport-wrapper framed --service Recommender --file ./thrift/recommender.thrift --api ./api/recommend.json --case case2 127.0.0.1:80
tenchmark -n 20000 -c 100 --service Recommender --file ./thrift/recommender.thrift --api ./api/recommend.json -c case1 127.0.0.1:80
```

### Thrift file
```thrift
namespace php Recommender.Thrift

struct MultiRequest {
1: required string user_id
       2: required string ad_id
       3: optional string city_name
       4: optional string first_category
       5: optional string second_category
       6: optional i32 size
}
service Recommender {

    string ping()
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        RecResponse fetchRecByMult(1:MultiRequest req)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        RecResponse fetchCtrByLR(1:string user_id, 2:string city, 3:i32 size)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)

        string fetchCategoryByUser(1:string user_id)
        throws (1: SystemException sys_exc,
                2: CodeException code_exc)
}

```
### API file
```json
{
    "case1": {
        "service": "Recommender",
        "function": "fetchRecByMult",
        "args": {
            "1": {
                "1": "etc_user_id",
                "2": "etc_ad_id",
                "3": "etc_city_name",
                "4": "etc_first_cate",
                "5": "etc_second_cate",
                "6": "etc_size"
            }
        }
    },
    "case2": {
        "service": "Recommender",
        "function": "fetchRecByMult",
        "args": {
            "1": {
                "1": "etc_user_id",
                "2": "etc_ad_id",
                "3": "etc_city_name",
                "4": "etc_first_cate",
                "5": "etc_second_cate",
                "6": "etc_size"
            }
        }
    },
}
```

