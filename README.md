# ali-logger-golang
aliyun k8s logger

## 安转
```shell script
go install github.com/yunkeCN/ali-logger-golang
```

## 使用
```go

import "github.com/yunkeCN/ali-logger-golang/logger"

// 在程序入口保证初始化成功
logger.Init(logger.Options{ProjectName: "middleman-server2", IsDev: true})

logger.Businessf("connect to http://localhost:%s/ for GraphQL playground", port)
logger.Business("connect to http://localhost:%s/ for GraphQL playground")
logger.Accessf("connect to http://localhost:%s/ for GraphQL playground", port)
logger.Access("connect to http://localhost:%s/ for GraphQL playground")
logger.Errorf("connect to http://localhost:%s/ for GraphQL playground", port)
logger.Error("connect to http://localhost:%s/ for GraphQL playground")

// 使用日志分组输出  WithTag
//增加、取消调用者输出  WithCaller
//增加、取消堆栈输出    WithStack 
// 增加通用字段输出     WithCommonField  WithCommonFields
// 增加自定义字段输出，统一放到attach节点     WithField  WithFields
// 取消自定义字段输出     ClearFields

log := logger.WithTag("feedback-online").
	WithStack(true).
	WithCommonField("requet.id","12345678").WithCommonFields(map[string]interface{}{"trace.id":"87654321"}).
	WithField("id","1").WithFields(map[string]interface{}{"name":"test","age":18})

log.Businessf("connect to http://localhost:%s/ for GraphQL playground", port)
log.Business("connect to http://localhost:%s/ for GraphQL playground")
log.Accessf("connect to http://localhost:%s/ for GraphQL playground", port)
log.Access("connect to http://localhost:%s/ for GraphQL playground")

log = logger.WithTag("feedback-online").
	WithCaller(true).
	WithCommonField("requet.id","12345678").WithCommonFields(map[string]interface{}{"trace.id":"87654321"}).
	WithField("id","1").WithFields(map[string]interface{}{"name":"test","age":18})
log.Errorf("connect to http://localhost:%s/ for GraphQL playground", port)
log.Error("connect to http://localhost:%s/ for GraphQL playground")
```

## 结合gin
```go
gin.DefaultWriter = logger.AccessWriter
gin.DefaultErrorWriter = logger.ErrorWriter

router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
    m := map[string]interface{}{
        "ClientIP":      param.ClientIP,
        "TimeStamp":     param.TimeStamp.Format(time.RFC1123),
        "Method":        param.Method,
        "Path":          param.Path,
        "Request.Proto": param.Request.Proto,
        "StatusCode":    param.StatusCode,
        "Latency":       param.Latency,
        "Agent":         param.Request.UserAgent(),
    }
    if param.ErrorMessage != "" {
        m["ErrorMessage"] = param.ErrorMessage
    }
    empData, err := json.Marshal(m)
    if err != nil {
        return ""
    }
    return string(empData) + "\n"
}))
```
