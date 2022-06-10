package controller

import (
	"github.com/gin-gonic/gin"
	"sshtunnelweb/app/myjwt"
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/util"
)

type Base struct{}

func (b *Base) Login(c *gin.Context) {
	resp := util.NewResult(c)
	type userinfo struct {
		Name   string `json:"username" form:"username"`
		Passwd string `json:"password" form:"password"`
	}
	var postInfo userinfo

	c.ShouldBind(&postInfo)
	res := myorm.User{}
	err := global.DB.Model(&myorm.User{}).Where(&myorm.User{Username: postInfo.Name, Passwd: postInfo.Passwd}).First(&res).Error

	if err == nil {

		getIDFromDB := res.ID
		getIPFromReq := util.GetRealIp(c)
		// fmt.Println(getIPFromReq)
		tmpFind := myorm.User{}
		err := global.DB.Model(&myorm.User{}).Where(&myorm.User{Ip: getIPFromReq}).Not("id = ?", getIDFromDB).First(&tmpFind).Error
		if err == nil {
			if tmpFind.ID != res.ID {
				global.Logger.Info(res.Username + "登录，ip被" + tmpFind.Username + "占用")
				resp.Error(500, tmpFind.Username+"-已占用登录IP")
				return
			}
		}
		if res.Ip != getIPFromReq {
			if err := global.DB.Model(&myorm.User{}).Where(&myorm.User{Username: postInfo.Name, Passwd: postInfo.Passwd}).Updates(&myorm.User{Ip: getIPFromReq}).Error; err != nil {
				global.Logger.Error(res.Username + "登录时更新登录IP失败：" + err.Error())
				resp.Error(500, "更新登录IP错误")
				return
			}
			for k, v := range global.LocalPortAndUserIP {
				if _, ok := v[res.Ip]; ok {
					delete(global.LocalPortAndUserIP[k], res.Ip)
					global.LocalPortAndUserIP[k][getIPFromReq] = "1"
				}
			}
		}
		j := myjwt.NewJWT()

		customClaims := j.CreateClaims(myjwt.BaseClaims{
			Username: postInfo.Name,
			ID:       getIDFromDB,
			IP:       getIPFromReq,
		})
		token, err := j.CreateToken(customClaims)
		if err != nil {
			global.Logger.Error(res.Username + "登录时，生产token失败：" + err.Error())
			resp.Error(500, err.Error())
			return
		}
		resp.Success(gin.H{
			"token": token,
			"info": gin.H{
				"username": postInfo.Name,
			},
			"ad": res.IsAdmin,
		})
		return
	} else {
		global.Logger.Error(err.Error())
		resp.Error(500, "登录信息错误")
		return
	}

}

func (b *Base) Register(c *gin.Context) {
	resp := util.NewResult(c)
	type userinfo struct {
		Name   string `json:"username"`
		Passwd string `json:"password"`
	}
	var postInfo userinfo

	c.ShouldBind(&postInfo)
	ip := util.GetRealIp(c)

	tmpFind := myorm.User{}
	if global.DB.Model(&myorm.User{}).Where(&myorm.User{Ip: ip}).First(&tmpFind).Error == nil {
		global.Logger.Error(postInfo.Name + ": " + tmpFind.Username + "-已占用你的IP")
		resp.Error(500, tmpFind.Username+"-已占用你的IP")
		return
	}
	if global.DB.Model(&myorm.User{}).Where(&myorm.User{Username: postInfo.Name}).First(&tmpFind).Error == nil {
		global.Logger.Error(postInfo.Name + ": " + tmpFind.Username + "-已占用")
		resp.Error(500, "用户名已被占用")
		return
	}
	if err := global.DB.Model(&myorm.User{}).Create(&myorm.User{Username: postInfo.Name, Passwd: postInfo.Passwd, Ip: ip}).Error; err != nil {
		global.Logger.Error(postInfo.Name + ": " + "创建账户失败: " + err.Error())
		resp.Error(500, "创建账户失败")
		return
	} else {
		resp.Success(gin.H{})
		return
	}
}
