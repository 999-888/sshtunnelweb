package controller

import (
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
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

	tmpdata := []myorm.User{}
	if global.DB.Preload("Conn").Find(&tmpdata).Error != nil {
		resp.Success(nil)
		return
	} else {
		data := []gin.H{}
		for _, k := range tmpdata {
			tmpconn := []gin.H{}
			if !k.IsAdmin {
				for _, c := range k.Conn {
					tmpconn = append(tmpconn, gin.H{
						"id":      c.ID,
						"svcname": c.Svcname,
					})
				}
			}
			data = append(data, gin.H{
				"id":       k.ID,
				"username": k.Username,
				"ip":       k.Ip,
				"isadmin":  k.IsAdmin,
				"conn":     tmpconn,
			})
		}
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
		global.Logger.Error("listuser: " + err.Error())
		resp.Error(500, "未获取到全部参数")
		return
	}
	delUser := myorm.User{}
	if global.DB.Model(&myorm.User{}).Where("id = ?", postInfo.ID).First(&delUser).Error != nil {
		global.Logger.Error(postInfo.ID + " 删除用户未找到")
		resp.Error(500, "删除用户未找到")
		return
	}
	if delUser.Username == global.CF.Admin.Name {
		global.Logger.Error("禁止删除初始化admin账户")
		resp.Error(500, "该用户不允许删除")
		return
	}
	if err := global.DB.Model(&delUser).Association("Conn").Clear(); err != nil {
		resp.Error(500, "删除授权服务关联失败")
		return
	}
	if err := global.DB.Model(&myorm.User{}).Delete(&myorm.User{}, postInfo.ID).Error; err != nil {
		global.Logger.Error(err.Error())
		resp.Error(500, "删除用户失败")
		return
	}
	resp.Success(nil)
	return
}

func UpdateUser(ctx *gin.Context) {
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
		global.Logger.Error(userinfo.Username + ": 更新用户信息：不是admin用户")
		resp.Error(500, "不是admin用户")
		return
	}
	type userInfo struct {
		Username string `form:"username" json:"username" binding:"required"`
		// Conn     []uint `form:"conn" json:"conn" binding:"required"`
		IsAdmin *bool `form:"isadmin" json:"isadmin" binding:"required"`
		// Ip      string `form:"ip" json:"ip" binding:"required"`
		Id uint `form:"id" json:"id" binding:"required"`
	}
	var postInfo userInfo
	if err := ctx.ShouldBind(&postInfo); err != nil {
		global.Logger.Error("update user； " + err.Error())
		resp.Error(500, "获取参数失败")
		return
	}
	// if err := global.DB.Model(&myorm.User{}).Where("id = ?", postInfo.Id).Updates(myorm.User{
	// Username: postInfo.Username,
	// IsAdmin:  *(postInfo.IsAdmin),
	// Ip:       postInfo.Ip,}
	if err := global.DB.Model(&myorm.User{}).Where("id = ?", postInfo.Id).Update("IsAdmin", *(postInfo.IsAdmin)).Error; err != nil {
		global.Logger.Error(postInfo.Username + "更新信息失败")
		resp.Error(500, "更新失败")
		return
	}
	// tmpuser := myorm.User{}
	// if err := global.DB.Model(&myorm.User{}).First(&tmpuser, postInfo.Id).Error; err != nil {
	// 	global.Logger.Error(err.Error())
	// 	resp.Error(500, err.Error())
	// 	return
	// }
	// 更新关联的服务，还需要判断服务是否打开了端口，未打开，还需要打开端口，不在这里更新，让用户自己申请
	// tmpconn := []myorm.Conn{}
	// if err := global.DB.Model(&myorm.Conn{}).Find(&tmpconn, postInfo.Conn).Error; err != nil {
	// 	global.Logger.Error(err.Error())
	// 	resp.Error(500, "更新关联的服务识别")
	// 	return
	// }
	// if err := global.DB.Model(&tmpuser).Association("Conn").Replace(tmpconn); err != nil {
	// 	global.Logger.Error(err.Error())
	// 	resp.Error(500, "更新关联的服务识别")
	// 	return
	// }
	resp.Success(nil)
	return
}
