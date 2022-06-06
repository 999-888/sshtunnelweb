package controller

import (
	"fmt"
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
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

	tmpdata := []resps.User{}
	if global.DB.Model(&myorm.User{}).Find(&tmpdata).Error != nil {
		resp.Success(nil)
		return
	} else {
		// fmt.Println(tmpdata)
		// item {"is": 1, "username": "root", "isadmin": true, "conn": [{"ID": 1, "svcname": },{"ID": 2, "svcname": }]}
		// return data = []item
		// data := []map[string]interface{}{}
		data := []gin.H{}
		for _, k := range tmpdata {
			tmpconn := []resps.RemoteSelect{}
			if err := global.DB.Model(&myorm.User{}).Where("id = ?", k.ID).Association("Conn").Find(&tmpconn); err != nil {
				global.Logger.Error("查找" + k.Username + "关联的conn失败: " + err.Error())
				resp.Error(500, "查找"+k.Username+"关联的conn失败: "+err.Error())
				return
			}
			// data = append(data, map[string]insterface{}{
			data = append(data, gin.H{
				"id":       k.ID,
				"username": k.Username,
				"ip":       k.Ip,
				"isadmin":  k.IsAdmin,
				"conn":     tmpconn,
			})
		}
		fmt.Println(data)
		resp.Success(data)
		return
	}

}

func DelUser(ctx *gin.Context) {
	resp := util.NewResult(ctx)
	userID, ok := ctx.Get("userid")
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
	type delData struct {
		ID string `form:"id" json:"id" binding:"required"`
	}
	var postInfo delData
	if err := ctx.ShouldBind(&postInfo); err != nil {
		fmt.Println(postInfo)
		resp.Error(500, "未获取到全部参数")
		return
	}
	delUser := myorm.User{}
	if global.DB.Model(&myorm.User{}).Where("id = ?", postInfo.ID).First(&delUser).Error != nil {
		resp.Error(500, "删除用户未找到")
		return
	}
	if delUser.Username == global.CF.Admin.Name {
		resp.Error(500, "该用户不允许删除")
		return
	}
	if global.DB.Model(&myorm.User{}).Where("id = ?", postInfo.ID).Association("Conn").Clear() != nil {
		resp.Error(500, "删除授权服务关联失败")
		return
	}
	if global.DB.Model(&myorm.User{}).Delete(&myorm.User{}, postInfo.ID).Error != nil {
		resp.Error(500, "删除用户失败")
		return
	}
	resp.Success(nil)
	return
}

func UpdateUser(ctx *gin.Context) {
	resp := util.NewResult(ctx)
	type sshinfo struct {
		Username string `form:"username" json:"username" binding:"required"`
		Passwd   string `form:"passwd" json:"passwd" binding:"required"`
		Port     string `form:"port" json:"port" binding:"required"`
		Host     string `form:"host" json:"host" binding:"required"`
		Id       int    `form:"id" json:"id" binding:"required"`
	}
	var postInfo sshinfo
	if err := ctx.ShouldBind(&postInfo); err != nil {
		fmt.Println(postInfo)
		resp.Error(500, "获取参数失败")
		return
	}
	if err := global.DB.Model(&myorm.Sshinfo{}).Where("id = ?", postInfo.Id).Updates(myorm.Sshinfo{
		Username: postInfo.Username,
		Passwd:   postInfo.Passwd,
		Port:     postInfo.Port,
		Host:     postInfo.Host,
	}).Error; err != nil {
		resp.Error(500, "更新失败")
		return
	}
	resp.Success(nil)
	return
}
