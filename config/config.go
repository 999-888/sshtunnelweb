package readconfig

import (
	// "fmt"
	// "github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init_configfile(configfile string, cf interface{}) error {

	viper.SetConfigFile(configfile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(cf); err != nil {
		return err
	}
	// 实时探测配置文件的变化
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	if err := viper.Unmarshal(cf); err != nil {
	// 		panic(err)
	// 	}
	// })
	return nil
}
