package controller

import (
	"fmt"
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

//展示 remote在使用的用户，用户为0时，才可删除

func ListSshtunnelRemote(c *gin.Context) {
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
		resp.Error(401, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}
	tmpres := []resps.Remote{}
	if global.DB.Model(&myorm.Conn{}).Find(&tmpres).Error != nil {
		resp.Success(nil)
		return
	} else {
		resp.Success(&tmpres)
		return

	}

}

func ListSshtunnelRemoteSelect(c *gin.Context) {
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
		resp.Error(403, "该用户未在db中查到")
		return
	}
	tmpres := []resps.RemoteSelect{}
	if global.DB.Model(&myorm.Conn{}).Find(&tmpres).Error != nil {
		resp.Success(nil)
		return
	} else {
		data := []map[string]interface{}{}
		for _, k := range tmpres {
			data = append(data, map[string]interface{}{
				"value": k.ID,
				"label": k.Svcname,
			})
		}
		resp.Success(&data)
		return

	}

}

func AddSshtunnelRemote(c *gin.Context) {
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
		resp.Error(403, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}
	type addRemote struct {
		Remote  string `form:"remote" json:"remote" binding:"required"`
		Svcname string `form:"svcname" json:"svcname" binding:"required"`
	}
	var postInfo addRemote

	if err := c.ShouldBind(&postInfo); err != nil {
		resp.Error(500, "获取参数失败")
		return
	}
	tmpres := myorm.Conn{
		Remote:  postInfo.Remote,
		Svcname: postInfo.Svcname,
	}
	if err := global.DB.Model(&myorm.Conn{}).Create(&tmpres).Error; err != nil {
		global.Logger.Error("创建" + postInfo.Remote + " - " + postInfo.Svcname + "失败: " + err.Error())
		resp.Error(500, "db 错误")
		return
	}
	resp.Success(nil)
	return
}

func UpdateSshtunnelRemote(c *gin.Context) {
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
		resp.Error(403, "该用户未在db中查到")
		return
	}
	if !userinfo.IsAdmin {
		global.Logger.Error("不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}
	type addRemote struct {
		Remote  string `form:"remote" json:"remote" binding:"required"`
		Svcname string `form:"svcname" json:"svcname" binding:"required"`
		ID      int    `form:"id" joson:"id" binding:"required"`
	}
	var postInfo addRemote

	if err := c.ShouldBind(&postInfo); err != nil {
		// fmt.Println(postInfo)
		resp.Error(500, "获取参数失败")
		return
	}
	tmpres := myorm.Conn{
		Remote:  postInfo.Remote,
		Svcname: postInfo.Svcname,
	}
	if err := global.DB.Model(&myorm.Conn{}).Where("id = ?", postInfo.ID).Updates(&tmpres).Error; err != nil {
		global.Logger.Error("更新" + postInfo.Remote + " - " + postInfo.Svcname + "失败: " + err.Error())
		resp.Error(500, "db 错误")
		return
	}
	resp.Success(nil)
	return
}

func DelSshtunnelRemote(c *gin.Context) {
	resp := util.NewResult(c)
	type delRemote struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}
	var postInfo delRemote

	if err := c.ShouldBind(&postInfo); err != nil {
		resp.Error(500, "获取参数失败")
		return
	}
	findConn := myorm.Conn{}
	if err := global.DB.Model(&myorm.Conn{}).First(&findConn, postInfo.ID); err != nil {
		resp.Error(500, "db err")
		return
	} else {
		if findConn.Local == "0" {
			if err := global.DB.Model(&myorm.Conn{}).Delete(&myorm.Conn{}, postInfo.ID).Error; err != nil {
				global.Logger.Error(fmt.Sprintf("删除conn %d 失败: %s", postInfo.ID, err.Error()))
				resp.Error(500, "db 错误")
				return
			}
		} else {
			resp.Error(401, "还有用户在使用")
			return
		}
	}
}
