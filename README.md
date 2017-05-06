# zlog
golang log library
## 简介
zlog是一个golang日志库, 轻量简单,支持日志分级, 分为debug, trace, info, error,fatal,优先级依次递增.
日志文件按照日期分割,格式为: 程序名.日期.log

## 示例
```go
package main

import (
	"github.com/wangzhen625/zlog"
)

func main() {

	//设定文件路径,不存在则自动创建,设定最小输出等级
	zlog.InitLogger("./log", zlog.DEBUG_LOG)

	zlog.Error("error test")
	//支持format格式
	zlog.Error("error test: %s", "zlog")
	zlog.Info("info test")
	zlog.Info("info test: %s", "zlog")
	zlog.Trace("trace test")
	zlog.Trace("trace test: %s", "zlog")
	zlog.Debug("debug test")
	zlog.Debug("debug test: %s", "zlog")
	// zlog.Fatal("fatal test")
	// zlog.Fatal("fatal test: %s", "zlog")
}
		
```

##输出样式
![ex](/image/ex.png)



