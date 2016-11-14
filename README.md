# zlog
golang log library
## 简介
	zlog是一个golang日志库, 功能十分简单,支持日志分级, 分为debug, notice, info, error四个等级,优先级依次递增.
	日志文件按照日期分割,格式为: 程序名.日期.log
## 示例
	···
		package main
		
		import (
			"github.com/wangzhen625/zlog"
		)

		func main() {

			//设定文件路径,不存在则自动创建,设定最小输出等级
			logger := log.InitLogger("../log", log.LEVEL_DEBUG)

			logger.Error("error test")
			//支持format格式
			logger.Error("error test: %s", "zlog")
			logger.Info("info test")
			logger.Info("info test: %s", "zlog")
			logger.Notice("notice test")
			logger.Notice("notice test: %s", "zlog")
			logger.Debug("debug test")
			logger.Debug("debug test: %s", "zlog")
		}

	...

##输出样式
	...
		2016-11-13 23:41:48 [ Error]: error test (E:/Go/goProjects/src/test/mylogtest.go:12) 
		2016-11-13 23:41:48 [ Error]: error test: zlog (E:/Go/goProjects/src/test/mylogtest.go:13) 
		2016-11-13 23:41:48 [ Info ]: info test (E:/Go/goProjects/src/test/mylogtest.go:14) 
		2016-11-13 23:41:48 [ Info ]: info test: zlog (E:/Go/goProjects/src/test/mylogtest.go:15) 
		2016-11-13 23:41:48 [Notice]: notice test (E:/Go/goProjects/src/test/mylogtest.go:16) 
		2016-11-13 23:41:48 [Notice]: notice test: zlog (E:/Go/goProjects/src/test/mylogtest.go:17) 
		2016-11-13 23:41:48 [ Debug]: debug test (E:/Go/goProjects/src/test/mylogtest.go:18) 
		2016-11-13 23:41:48 [ Debug]: debug test: zlog (E:/Go/goProjects/src/test/mylogtest.go:19) 
	...



