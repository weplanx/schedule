# Schedule gRPC

Manage scheduled tasks using gRPC

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainonly/schedule-microservice?style=flat-square)](https://github.com/kainonly/schedule-microservice)
[![Github Actions](https://img.shields.io/github/workflow/status/kainonly/schedule-microservice/release?style=flat-square)](https://github.com/kainonly/schedule-microservice/actions)
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
```

## Configuration

For configuration, please refer to `config/config.example.yml`

- **debug** `string` Start debugging, ie `net/http/pprof`, access address is `http://localhost:6060`
- **listen** `string` Microservice listening address
- **logging** `object` Log configuration
    - **storage** `string` Local log storage directory
    - **transfer** `object` [elastic-transfer](https://github.com/kainonly/elastic-transfer) service
      - **listen** `string` host
      - **id** `string` transfer id

## Service

The service is based on gRPC and you can view `router/router.proto`

```proto
syntax = "proto3";
package schedule;
service Router {
    rpc Get (GetParameter) returns (GetResponse) {
    }

    rpc Lists (ListsParameter) returns (ListsResponse) {
    }

    rpc All (NoParameter) returns (AllResponse) {
    }

    rpc Put (PutParameter) returns (Response) {
    }

    rpc Delete (DeleteParameter) returns (Response) {
    }

    rpc Running (RunningParameter) returns (Response) {
    }
}

message NoParameter {
}

message Response {
    uint32 error = 1;
    string msg = 2;
}

message EntryOption {
    string cron_time = 1;
    string url = 2;
    bytes headers = 3;
    bytes body = 4;
}

message EntryOptionWithTime {
    string cron_time = 1;
    string url = 2;
    bytes headers = 3;
    bytes body = 4;
    int64 next_date = 5;
    int64 last_date = 6;
}

message Information {
    string identity = 1;
    bool start = 2;
    string time_zone = 3;
    map<string, EntryOptionWithTime> entries = 4;
}

message GetParameter {
    string identity = 1;
}

message GetResponse {
    uint32 error = 1;
    string msg = 2;
    Information data = 3;
}

message ListsParameter {
    repeated string identity = 1;
}

message ListsResponse {
    uint32 error = 1;
    string msg = 2;
    repeated Information data = 3;
}

message AllResponse {
    uint32 error = 1;
    string msg = 2;
    repeated string data = 3;
}

message PutParameter {
    string identity = 1;
    string time_zone = 2;
    bool start = 3;
    map<string, EntryOption> entries = 4;
}

message DeleteParameter {
    string identity = 1;
}

message RunningParameter {
    string identity = 1;
    bool running = 2;
}
```

#### rpc Get (GetParameter) returns (GetResponse) {}

Get job information

- GetParameter
  - **identity** `string` job id
- GetResponse
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback
  - **data** `Information` result
    - **identity** `string` job id
    - **start** `bool` operating status
    - **time_zone** `string` time zone, [tz database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
    - **entries** `map<string, EntryOptionWithTime>` The task workflow
      - **string** task id
      - EntryOptionWithTime
        - **cron_time** `string` CronTab rule
        - **url** `string` callback hook url
        - **headers** `bytes`
        - **body** `bytes` 
        - **next_date** `int64` next run unixtime
        - **last_date** `int64` last run unixtime

```golang
client := pb.NewRouterClient(conn)
response, err := client.Get(
    context.Background(),
    &pb.GetParameter{
        Identity: "test",
    },
)
```

#### rpc Lists (ListsParameter) returns (ListsResponse) {}

Get job information in batches

- ListsParameter
  - **identity** `[]string` job IDs
- ListsResponse
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback
  - **data** `[]Information` result
    - **identity** `string` job id
    - **start** `bool` operating status
    - **time_zone** `string` time zone, [tz database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
    - **entries** `map<string, EntryOptionWithTime>` The task workflow
      - **string** task id
      - EntryOptionWithTime
        - **cron_time** `string` CronTab rule
        - **url** `string` callback hook url
        - **headers** `bytes`
        - **body** `bytes` 
        - **next_date** `int64` next run unixtime
        - **last_date** `int64` last run unixtime

```golang
client := pb.NewRouterClient(conn)
response, err := client.Lists(
    context.Background(),
    &pb.ListsParameter{
        Identity: []string{"test", "other"},
    },
)
```

#### rpc All (NoParameter) returns (AllResponse) {}

Get all job IDs

- NoParameter
- AllResponse
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback
  - **data** `[]string` job IDs

```golang
client := pb.NewRouterClient(conn)
response, err := client.All(
    context.Background(),
    &pb.NoParameter{},
)
```

#### rpc Put (PutParameter) returns (Response) {}

Add or update job

- PutParameter
  - **identity** `string` job id
  - **time_zone** `string` time zone, [tz database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
  - **start** `bool` operating status
  - **entries** `map<string, EntryOption>` The task workflow
    - **string** task id
    - EntryOptionWithTime
      - **cron_time** `string` CronTab rule
      - **url** `string` callback hook url
      - **headers** `bytes`
      - **body** `bytes` 
- Response
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback

```golang
client := pb.NewRouterClient(conn)
response, err := client.Put(
    context.Background(),
    &pb.PutParameter{
        Identity: "test",
        TimeZone: "Asia/Shanghai",
        Start:    true,
        Entries: map[string]*pb.EntryOption{
            "task1": &pb.EntryOption{
                CronTime: "*/5 * * * * *",
                Url:      "http://localhost:3000",
                Headers:  []byte(`{"x-token":"abc"}`),
                Body:     []byte(`{"name":"task1"}`),
            },
            "task2": &pb.EntryOption{
                CronTime: "*/10 * * * * *",
                Url:      "http://localhost:3000",
                Headers:  []byte(`{"x-token":"abc"}`),
                Body:     []byte(`{"name":"task2"}`),
            },
        },
    },
)
```

#### rpc Delete (DeleteParameter) returns (Response) {}

remvoe job

- DeleteParameter
  - **identity** `string` job id
- Response
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback

```golang
client := pb.NewRouterClient(conn)
response, err := client.Delete(
    context.Background(),
    &pb.DeleteParameter{
        Identity: "test",
    },
)
```

#### rpc Running (RunningParameter) returns (Response) {}

Change operation status

- RunningParameter
  - **identity** `string` job id
  - **running** `bool` operating status
- Response
  - **error** `uint32` error code, `0` is normal
  - **msg** `string` error feedback

```golang
client := pb.NewRouterClient(conn)
response, err := client.Running(
    context.Background(),
    &pb.RunningParameter{
        Identity: "test",
        Running:  false,
    },
)
```