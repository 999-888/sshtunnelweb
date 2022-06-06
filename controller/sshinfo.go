package controller

import (
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

func ListSshinfo(c *gin.Context) {
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

	tmpdata := []resps.Sshinfo{}
	if global.DB.Model(&myorm.Sshinfo{}).Find(&tmpdata).Error != nil {
		resp.Success(nil)
		return
	} else {
		// fmt.Println(tmpdata)
		resp.Success(tmpdata)
		return
	}

}

func AddSshinfo(ctx *gin.Context) {
	resp := util.NewResult(ctx)
	type addsshinfo struct {
		Username string `form:"Username" json:"Username" binding:"required"`
		Passwd   string `form:"Passwd" json:"Passwd" binding:"required"`
		Port     string `form:"Port" json:"Port" binding:"required"`
		Host     string `form:"Host" json:"Host" binding:"required"`
	}
	var postInfo addsshinfo
	if err := ctx.ShouldBind(&postInfo); err != nil {
		// fmt.Println(postInfo)
		resp.Error(500, "未获取到全部参数")
		return
	}
	// tmpdata := map[string]interface{}{}
	tmpdata := myorm.Sshinfo{}
	// if global.DB.Table("sshinfo").Take(&tmpdata).RowsAffected == 0 {
	if global.DB.Model(&myorm.Sshinfo{}).First(&tmpdata).Error != nil {
		// if err := global.DB.Table("sshinfo").Create(map[string]interface{}{
		if err := global.DB.Model(&myorm.Sshinfo{}).Create(&myorm.Sshinfo{
			Username: postInfo.Username,
			Passwd:   postInfo.Passwd,
			Port:     postInfo.Port,
			Host:     postInfo.Host,
		}).Error; err != nil {
			resp.Error(500, "新增失败")
			return
		}
		resp.Success(nil)
		return
	} else {
		resp.Error(500, "只允许一个ssh存在,请更新")
		return
	}
}

func UpdateSshinfo(ctx *gin.Context) {
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
		// fmt.Println(postInfo)
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
