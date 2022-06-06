package global

import (
	"fmt"
	"os"
	"sshtunnelweb/zaplog"

	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
)

//创建logger
func SetupLogger() {
	var err error
	if _, err := os.Stat(CF.Logs.Path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(CF.Logs.Path, os.ModePerm); err != nil {
				fmt.Println("创建日志目录失败")
				os.Exit(-1)
			}
		} else {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	filepath := CF.Logs.Path
	infofilename := fmt.Sprintf("%s-info", CF.Logs.Logfilename)
	warnfilename := fmt.Sprintf("%s-warn", CF.Logs.Logfilename)
	fileext := CF.Logs.Logfileext

	//infofilename,warnfilename,fileext string
	Logger, err = zaplog.GetInitLogger(filepath, infofilename, warnfilename, fileext)

	if err != nil {
		fmt.Println("创建日志失败", err)
		os.Exit(-1)
	}
	defer Logger.Sync()
}
