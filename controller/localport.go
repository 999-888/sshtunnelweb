package controller

import (
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

func ListLocalPort(c *gin.Context) {
	resp := util.NewResult(c)
	userID, ok := c.Get("userid")
	if !ok {
		global.Logger.Error("没有获取到jwt信息")
		resp.Error(500, "没有获取到jwt信息")
		return
	}
	userinfo := myorm.User{}
	if global.DB.Model(&myorm.User{}).First(&userinfo, userID).Error != nil {
		global.Logger.Error("该用户未在db中查到")
		resp.Error(500, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}
	tmpLocalPort := []resps.Conn{}

	if err := global.DB.Model(&myorm.Conn{}).Not("local = ?", "").Find(&tmpLocalPort).Error; err != nil {
		global.Logger.Error("查找错误: " + err.Error())
		resp.Success(nil)
		return
	}

	resp.Success(&tmpLocalPort)
	return
}

func ListOneUserLocalPort(c *gin.Context) {
	resp := util.NewResult(c)
	userID, ok := c.Get("userid")
	if !ok {
		global.Logger.Error("没有获取到jwt信息")
		resp.Error(500, "没有获取到jwt信息")
		return
	}
	userinfo := myorm.User{}
	if global.DB.Model(&myorm.User{}).First(&userinfo, userID).Error != nil {
		global.Logger.Error("该用户未在db中查到")
		resp.Error(500, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}

	type postdata struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}
	var postInfo postdata
	if err := c.ShouldBind(&postInfo); err != nil {
		resp.Error(500, "未获取到全部参数")
		return
	}
	tmpLocalPort := myorm.Conn{}

	if err := global.DB.Model(&myorm.Conn{}).Preload("User").First(&tmpLocalPort, postInfo.ID).Error; err != nil {
		global.Logger.Error("查找错误: " + err.Error())
		resp.Success(nil)
		return
	}
	tmpuser := []gin.H{}
	for _, k := range tmpLocalPort.User {
		tmpuser = append(tmpuser, gin.H{
			"id":       k.ID,
			"username": k.Username,
			"ip":       k.Ip,
		})
	}
	resp.Success(&tmpuser)
	return
}

func DelOneUserLocalPort(c *gin.Context) {
	resp := util.NewResult(c)
	userID, ok := c.Get("userid")
	if !ok {
		global.Logger.Error("没有获取到jwt信息")
		resp.Error(500, "没有获取到jwt信息")
		return
	}
	userinfo := myorm.User{}
	if global.DB.Model(&myorm.User{}).First(&userinfo, userID).Error != nil {
		global.Logger.Error("该用户未在db中查到")
		resp.Error(500, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}

	type postdata struct {
		ID     uint `form:"id" json:"id" binding:"required"`
		ConnID uint `form:"conn" json:"conn" binding:"required"`
	}
	var postInfo postdata
	if err := c.ShouldBind(&postInfo); err != nil {
		resp.Error(500, "未获取到全部参数")
		return
	}
	tmpUser := myorm.User{}

	if err := global.DB.Model(&myorm.User{}).First(&tmpUser, postInfo.ID).Error; err != nil {
		global.Logger.Error("查找错误: " + err.Error())
		resp.Error(500, err.Error())
		return
	}
	tmpConn := myorm.Conn{}
	if err := global.DB.Model(&myorm.Conn{}).First(&tmpConn, postInfo.ConnID).Error; err != nil {
		global.Logger.Error("查找错误: " + err.Error())
		resp.Error(500, err.Error())
		return
	}
	if err := global.DB.Model(&tmpUser).Association("Conn").Delete(&tmpConn); err != nil {
		global.Logger.Error(tmpUser.Username + "去关联 " + tmpConn.Svcname + "失败；" + err.Error())
		resp.Error(500, "取消关联失败")
		return
	}
	global.Logger.Info(tmpUser.Username + "被admin权限账户" + userinfo.Username + "去关联 " + tmpConn.Svcname + "成功")
	delete(global.LocalPortAndUserIP[tmpConn.Local], tmpUser.Ip)
	resp.Success(nil)
	return
}
