package global

// http 请求的日志，替换gin自带的log

import (
	// "fmt"
	"sshtunnelweb/zaplog"

	"go.uber.org/zap"
)

var (
	AccessLogger *zap.SugaredLogger
)

//创建logger
func SetupAccessLogger() error {
	var err error
	filepath := CF.AccessLogs.Path
	filename := CF.AccessLogs.Logfilename
	//warnfilename:= LogSetting.LogWarnFileName
	fileext := CF.AccessLogs.Logfileext
	// fmt.Println(filepath, filename, fileext)
	//infofilename,warnfilename,fileext string
	AccessLogger, err = zaplog.GetInitAccessLogger(filepath, filename, fileext)

	if err != nil {
		return err
	}
	defer AccessLogger.Sync()
	return nil
}
