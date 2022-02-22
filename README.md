# Weplanx Schedule

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/weplanx/schedule?style=flat-square)](https://github.com/weplanx/schedule)
[![Go Report Card](https://goreportcard.com/badge/github.com/weplanx/schedule?style=flat-square)](https://goreportcard.com/report/github.com/weplanx/schedule)
[![Release](https://img.shields.io/github/v/release/weplanx/schedule.svg?style=flat-square)](https://github.com/weplanx/schedule)
[![GitHub license](https://img.shields.io/github/license/weplanx/schedule?style=flat-square)](https://raw.githubusercontent.com/weplanx/schedule/main/LICENSE)

定时调度器，协助应用定时触发需要的任务

> 项目将以新的方式重新开发配套 weplanx ，新版本将以 `v*.*.*` 形式发布

## 客户端

在 `go.mod` 项目中

```shell
go get github.com/weplanx/schedule
```

简单使用

```go
package main

import (
	"context"
	"fmt"
	"github.com/weplanx/schedule/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 或者使用 TLS
	// certFile := "..."
	// creds, err := credentials.NewClientTLSFromFile(certFile, "")
	// if err != nil {
	// 	panic(err)
	// }
	// opts = append(opts, grpc.WithTransportCredentials(creds))

	schedule, err := client.New("127.0.0.1:6000", opts...)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	result, err := schedule.Get(ctx, []string{})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
```

## 部署服务

镜像源主要有：

- ghcr.io/weplanx/schedule:latest
- ccr.ccs.tencentyun.com/weplanx/schedule:latest（国内）

案例将使用 Kubernetes 部署编排，复制部署内容（需要根据情况做修改）：

1. 设置配置

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: schedule.cfg
data:
  config.yml: |
    address: ":6000"
    tls: <TLS配置，非必须>
      cert:
      key:
    node: <节点标识>
    database:
      uri: mongodb://<username>:<password>@<host>:<port>/<database>?authSource=<authSource>
      name: <数据库名>
      collection: <默认集合>
    transfer:
      address: "transfer-svc:6000"
      topic: schedule
      tls:
        cert:
```

2. 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: schedule
  name: schedule-deploy
spec:
  selector:
    matchLabels:
      app: schedule
  template:
    metadata:
      labels:
        app: schedule
    spec:
      containers:
        - image: ccr.ccs.tencentyun.com/weplanx/schedule:latest
          imagePullPolicy: Always
          name: schedule
          ports:
            - containerPort: 6000
          volumeMounts:
            - name: config
              mountPath: "/app/config"
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: schedule.cfg
            items:
              - key: "config.yml"
                path: "config.yml"
```

3. 设置入口，服务网关推荐采用 traefik 做更多处理

```yaml
apiVersion: v1
kind: Service
metadata:
  name: schedule-svc
spec:
  ports:
    - port: 6000
      protocol: TCP
  selector:
    app: schedule
```

## 滚动更新

复制模板内容，并需要自行定制触发条件，原理是每次patch将模板中 `${tag}` 替换为版本执行

```yml
spec:
  template:
    spec:
      containers:
        - image: ccr.ccs.tencentyun.com/weplanx/schedule:${tag}
          name: schedule
```

例如：在 Github Actions
中 `patch deployment schedule-deploy --patch "$(sed "s/\${tag}/${{steps.meta.outputs.version}}/" < ./config/patch.yml)"`，国内可使用**Coding持续部署**或**云效流水线**等。

## License

[BSD-3-Clause License](https://github.com/weplanx/schedule/blob/main/LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fweplanx%2Fschedule.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fweplanx%2Fschedule?ref=badge_large)