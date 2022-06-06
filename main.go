package main

import (
	// "fmt"
	"fmt"
	"sshtunnelweb/global"
	"sshtunnelweb/router"
)

func main() {
	// 读取启动时指定的配置文件
	global.ReadConfigFile()
	// 配置日志
	global.SetupLogger()
	global.SetupAccessLogger()
	// 初始化sqlite
	global.InitSqlite()

	r := router.Router()

	r.Run(fmt.Sprintf(":%s", global.CF.Run.Port))
}
