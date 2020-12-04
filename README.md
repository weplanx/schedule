# Schedule Microservice

Manage scheduled tasks using gRPC

[![Github Actions](https://img.shields.io/github/workflow/status/kain-lab/schedule-microservice/release?style=flat-square)](https://github.com/kain-lab/schedule-microservice/actions)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kain-lab/schedule-microservice?style=flat-square)](https://github.com/kain-lab/schedule-microservice)
[![Image Size](https://img.shields.io/docker/image-size/kainonly/schedule-microservice?style=flat-square)](https://hub.docker.com/r/kainonly/schedule-microservice)
[![Docker Pulls](https://img.shields.io/docker/pulls/kainonly/schedule-microservice.svg?style=flat-square)](https://hub.docker.com/r/kainonly/schedule-microservice)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/kainonly/schedule-microservice/master/LICENSE)

## Setup

Example using docker compose

```yaml
version: "3.8"
services: 
  schedule:
    image: kainonly/schedule-microservice
    restart: always
    volumes: 
      - ./schedule/config:/app/config
      - ./schedule/log:/app/log
    ports:
      - 6000:6000
      - 8080:8080
```

## Configuration

For configuration, please refer to `config/config.example.yml`

- **debug** `string` Start debugging, ie `net/http/pprof`, access address is `http://localhost:6060`
- **listen** `string` grpc server listening address
- **gateway** `string` API gateway server listening address
- **filelog** `string` file log
- **transfer** `object` [elastic-transfer](https://github.com/kain-lab/elastic-transfer) service
  - **listen** `string` host
  - **pipe** `string` id

## Service

The service is based on gRPC and you can view `router/router.proto`

```proto
syntax = "proto3";
package schedule;
option go_package = "schedule-microservice/gen/go/schedule";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";

service API {
  rpc Get (ID) returns (Option) {
    option (google.api.http) = {
      get: "/schedule",
    };
  }
  rpc Lists (IDs) returns (Options) {
    option (google.api.http) = {
      post: "/schedules",
      body: "*"
    };
  }
  rpc All (google.protobuf.Empty) returns (IDs) {
    option (google.api.http) = {
      get: "/schedules",
    };
  }
  rpc Put (Option) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      put: "/schedule",
      body: "*",
    };
  }
  rpc Delete (ID) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/schedule",
    };
  }
  rpc Running (Status) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      put: "/running",
      body: "*",
    };
  }
}

message ID {
  string id = 1;
}

message IDs {
  repeated string ids = 1;
}

message Option {
  string id = 1;
  string time_zone = 2;
  bool start = 3;
  map<string, Entry> entries = 4;
}

message Options {
  repeated Option options = 1;
}

message Entry {
  string cron_time = 1;
  string url = 2;
  bytes headers = 3;
  bytes body = 4;
  int64 next_date = 5;
  int64 last_date = 6;
}

message Status {
  string id = 1;
  bool running = 2;
}
```

## Get (ID) returns (Option)

Get job configuration

### RPC

- **ID**
  - **id** `string` job id
- **Option**
  - **id** `string` job id
  - **time_zone** `string` time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
  - **start** `bool` operating status
  - **entries** `map<string, Entry>` the task workflow
    - **string** task id
    - **Entry**
      - **cron_time** `string` cronTab rule
      - **url** `string` callback hook url
      - **headers** `bytes`
      - **body** `bytes` 
      - **next_date** `int64` next run unixtime
      - **last_date** `int64` last run unixtime

```golang
client := pb.NewAPIClient(conn)
response, err := client.Get(context.Background(), &pb.ID{
  Id: "debug-A",
})
```

### API Gateway

- **GET** `/schedule`

```http
GET /schedule?id=debug-C HTTP/1.1
Host: localhost:8080
```

## Lists (IDs) returns (Options)

Lists job configuration

### RPC

- **IDs**
  - **ids** `[]string` job IDs
- **Options**
  - **options** `[]Option` result
    - **id** `string` job id
    - **time_zone** `string` time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
    - **start** `bool` operating status
    - **entries** `map<string, Entry>` the task workflow
      - **string** task id
      - **Entry**
        - **cron_time** `string` cronTab rule
        - **url** `string` callback hook url
        - **headers** `bytes`
        - **body** `bytes` 
        - **next_date** `int64` next run unixtime
        - **last_date** `int64` last run unixtime

```golang
client := pb.NewAPIClient(conn)
response, err := client.Lists(context.Background(), &pb.IDs{
  Ids: []string{"debug-A"},
})
```

### API Gateway

- **POST** `/schedules`

```http
POST /schedules HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "ids":["debug-C"]
}
```

## All (google.protobuf.Empty) returns (IDs)

Get all job configuration identifiers

### RPC

- **IDs**
  - **ids** `[]string` job IDs

```golang
client := pb.NewAPIClient(conn)
response, err := client.All(context.Background(), &empty.Empty{})
```

### API Gateway

- **GET** `/schedules`

```http
GET /schedules HTTP/1.1
Host: localhost:8080
```

## Put (Option) returns (google.protobuf.Empty)

Put job configuration

### RPC

- **Option**
  - **id** `string` job id
  - **time_zone** `string` time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
  - **start** `bool` operating status
  - **entries** `map<string, Entry>` the task workflow
    - **string** task id
    - **Entry**
      - **cron_time** `string` cronTab rule
      - **url** `string` callback hook url
      - **headers** `bytes`
      - **body** `bytes` 

```golang
client := pb.NewAPIClient(conn)
_, err := client.Put(context.Background(), &pb.Option{
  Id:       "debug-A",
  TimeZone: "Asia/Shanghai",
  Start:    true,
  Entries: map[string]*pb.Entry{
    "entry-1": {
      CronTime: "*/10 * * * * *",
      Url:      "http://mac:3000/entry-1",
      Headers:  []byte(`{"x-token":"l51aM51gp43606o2"}`),
      Body:     []byte(`{"msg":"hello entry-A1"}`),
    },
    "entry-2": {
      CronTime: "*/20 * * * * *",
      Url:      "http://mac:3000/entry-2",
      Headers:  []byte(`{"x-token":"GGlxNXfMyJb5IKuL"}`),
      Body:     []byte(`{"msg":"hello entry-A2"}`),
    },
  },
})
```

### API Gateway

- **PUT** `/schedule`

```http
PUT /schedule HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "id": "debug-C",
    "time_zone": "Asia/Shanghai",
    "start": true,
    "entries": {
        "entry": {
            "cron_time": "*/10 * * * * *",
            "url": "http://mac:3000/entry",
            "headers": "eyJ4LXRva2VuIjoiMTIzNDU2In0=",
            "body": "eyJtc2ciOiJoZWxsbyBlbnRyeSJ9"
        }
    }
}
```

## Delete (ID) returns (google.protobuf.Empty)

Remove job configuration

### RPC

- **ID**
  - **id** `string` job id

```golang
client := pb.NewAPIClient(conn)
_, err := client.Delete(context.Background(), &pb.ID{
  Id: "debug-A",
})
```

### API Gateway

- **DELETE** `/schedule`

```http
DELETE /schedule?id=debug-C HTTP/1.1
Host: localhost:8080
```

## Running (Status) returns (google.protobuf.Empty)

Change job status

### RPC

- **Status**
  - **id** `string` job id
  - **running** `bool` operating status

```golang
client := pb.NewAPIClient(conn)
_, err := client.Running(context.Background(), &pb.Status{
  Id:      "debug-A",
  Running: false,
})
```

### API Gateway

- **PUT** `/running`

```http
PUT /running HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "id": "debug-C",
    "running": false
}
```