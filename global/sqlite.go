package global

import (
	"fmt"
	"os"
	"path/filepath"
	"sshtunnelweb/myorm"
	"sshtunnelweb/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	// dbpath string = CF.Sqlite.Path
	// dbfile string = CF.Sqlite.Filename
)

func InitSqlite() {
	// fmt.Println(CF.Sqlite.Path)
	if _, err := os.Stat(CF.Sqlite.Path); err != nil {
		fmt.Println(err.Error())
		if os.IsNotExist(err) {
			if err := os.MkdirAll(CF.Sqlite.Path, os.ModePerm); err != nil {
				Logger.Error("创建db日志目录失败")
				fmt.Println("创建db日志目录失败")
				os.Exit(-1)
			}
		} else {
			Logger.Error(err.Error())
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}
	dbpf := filepath.Join(CF.Sqlite.Path, CF.Sqlite.Filename)
	db, err := gorm.Open(sqlite.Open(dbpf), &gorm.Config{})
	if err != nil {
		Logger.Error("连接db失败")
		fmt.Println("连接db失败")
		os.Exit(-1)
	}

	DB = db
	// 迁移 schema
	for _, k := range myorm.GetAll() {
		db.AutoMigrate(k)
	}

	// 创建超级账户
	adminUser := myorm.User{Username: CF.Admin.Name, Passwd: CF.Admin.Passwd, IsAdmin: true}
	resUser := myorm.User{}
	if db.Model(&myorm.User{}).Where(&adminUser).First(&resUser).RowsAffected == 0 {
		if err := db.Model(&myorm.User{}).Create(&adminUser).Error; err != nil {
			Logger.Error("启动ing：创建admin账户出错")
			fmt.Println("启动ing：创建admin账户出错")
			os.Exit(-1)
		}
	}
	// 程序重启，需要重新创建那些已经打开的 本地端口 和 交互数据的监听协程
	findConn := []myorm.Conn{}
	if err := db.Model(&myorm.Conn{}).Not("local = ?", "").Find(&findConn).Error; err != nil {
		// if err.Error()
		Logger.Error("启动ing：查询db出错" + err.Error())
		fmt.Println("启动ing：查询db出错" + err.Error())
		os.Exit(-1)
	} else {
		for _, k := range findConn {
			if ST == nil {
				if err := StartST(); err != nil {
					Logger.Error("启动ing：连接远程机出错")
					fmt.Println("启动ing：连接远程机出错")
					os.Exit(-1)
				}
			}
			local := util.GetOnePort()
			local_listen, err := StartLocalListen(":" + k.Local)
			if err != nil {
				Logger.Error("启动ing：" + k.Local + ": 本地监听失败")
				fmt.Println("启动ing：" + k.Local + ": 本地监听失败")
				os.Exit(-1)
			}
			// p := []*net.Listener{}
			// p = append(p, &local_listen)
			GlobalSshtunnelInfo[local] = &local_listen
			go StartTunnel(local_listen, k.Remote, ST)
		}
	}
}
