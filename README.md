# zlog
golang log library
## 简介
zlog是一个golang日志库, 轻量简单,具有以下特点：

1.日志分级，分为Debug、Trace、Info、Error、Fatal五个等级，优先级依次递增。

2.文件分割

 -  按大小分割，默认50M左右分割一个文件。
 -  按日期分割，每天零点生成当天的日志文件。 


文件分割两种方式混合使用，生成文件名格式：serv.20170517-143344.slave01.12439.log

意义分别为：程序名.日期-时间.主机名.进程号.log


## 示例

```go
package main

import (
	"github.com/wangzhen625/zlog"
)

func main() {

	//设定文件路径,不存在则自动创建,设定最小输出等级
	zlog.InitLogger("./log", zlog.DebugLevel)

	zlog.Error("error test")
	//支持format格式
	zlog.Error("error test: %s", "zlog")
	zlog.Info("info test")
	zlog.Info("info test: %s", "zlog")
	zlog.Trace("trace test")
	zlog.Trace("trace test: %s", "zlog")
	zlog.Debug("debug test")
	zlog.Debug("debug test: %s", "zlog")
	zlog.Fatal("fatal test")
	zlog.Fatal("fatal test: %s", "zlog")
}
		
```

##输出样式
![tmp](/image/tmp.png)



