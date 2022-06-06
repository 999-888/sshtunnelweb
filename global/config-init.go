package global

import (
	"flag"
	"fmt"
	"os"
	sshtunnelweb "sshtunnelweb/config"
)

type Sqlite struct {
	Path     string `yaml:"path"`
	Filename string `yamml:"filename"`
}

type Logs struct {
	Path        string `yaml:"path"`
	Logfilename string `yaml:"logfilename"`
	Logfileext  string `yaml:"logfileext"`
}

type AccessLogs struct {
	Path        string `yaml:"path"`
	Logfilename string `yaml:"logfilename"`
	Logfileext  string `yaml:"logfileext"`
}

type Run struct {
	Port string `yaml:"port"`
}

type Jwt struct {
	ExpiresTime string `yaml:"expires_time" mapstructure:"expires_time"`
	BufferTime  string `yaml:"buffer_time" mapstructure:"buffer_time"`
	Issuer      string `yaml:"issuer"`
}

type Admin struct {
	Name   string `yaml:"name"`
	Passwd string `yaml:"passwd"`
}

type Config struct {
	Sqlite     `yaml:"sqlite"`
	Logs       `yaml:"logs"`
	AccessLogs `yaml:"accesslogs"`
	Run        `yaml:"run"`
	Jwt        `yaml:"jwt"`
	Admin      `yaml:"admin"`
}

var CF = new(Config)

func ReadConfigFile() {
	var config_file string
	flag.StringVar(&config_file, "c", "./config.yaml", "配置文件路径")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `filerserver -H
		Usage: fileserver [-c 指定配置文件]`)
		flag.PrintDefaults()
	}
	flag.Parse()
	if _, err := os.Stat(config_file); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(config_file, "  文件不存在")
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}

	if err := sshtunnelweb.Init_configfile(config_file, CF); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
