package controller

import (
	"fmt"
	"sshtunnelweb/global"
	"syscall"

	"github.com/gin-gonic/gin"
)

func TT(c *gin.Context) {
	fmt.Println("tt")
	global.Logger.Info("zy zy zy")
	global.Logger.Errorf("%s-err: %s", "zdy", "666")
	fmt.Println(syscall.Getpid())
	c.JSON(200, gin.H{
		"msg": global.CF,
	})
	return
}
