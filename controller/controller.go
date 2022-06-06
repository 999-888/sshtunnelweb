package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sshtunnelweb/global"
)

func TT(c *gin.Context) {
	fmt.Println("tt")
	global.Logger.Info("zy zy zy")
	global.Logger.Errorf("%s-err: %s", "zdy", "666")
	c.JSON(200, gin.H{
		"msg": global.CF,
	})
	return
}
