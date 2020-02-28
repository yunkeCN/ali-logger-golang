# ali-logger-golang
aliyun k8s logger

## 安转
```shell script
go install github.com/yunkeCN/ali-logger-golang
```

## 使用
```go

import "github.com/yunkeCN/ali-logger-golang/logger"

logger.Init(logger.Options{ProjectName: "middleman-server2", IsDev: true})

logger.Businessf("connect to http://localhost:%s/ for GraphQL playground", port)
logger.Accessf("connect to http://localhost:%s/ for GraphQL playground", port)
logger.Errorf("connect to http://localhost:%s/ for GraphQL playground", port)
```
