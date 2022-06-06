package controller

import (
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

// add del  post 获取到的参数
type sshTunnel struct {
	// Local  string `form:"local" json:"local" binding:"required"`
	// Remote string `form:"remote" json:"remote" binding:"required"`
	// 由于选择器的存在，svcname的值变为了id的值
	ID uint `form:"svcname" json:"svcname" binding:"required"`
}

func ListSshtunnel(c *gin.Context) {
	resp := util.NewResult(c)
	userID, ok := c.Get("userid")
	if !ok {
		resp.Error(500, "没有获取到jwt信息")
		return
	}
	userinfo := myorm.User{}
	if global.DB.Model(&myorm.User{}).First(&userinfo, userID).Error != nil {
		resp.Error(500, "该用户未在db中查到")
		return
	}
	tmpres := []resps.Conn{}
	var err error
	if userinfo.IsAdmin {
		err = global.DB.Model(&myorm.Conn{}).Not("local = ?", "").Find(&tmpres).Error
	} else {
		tmp := myorm.User{}
		err = global.DB.Preload("Conn").First(&tmp, userID).Error
		for _, k := range tmp.Conn {
			tmpres = append(tmpres, resps.Conn{ID: k.ID, Local: k.Local, Svcname: k.Svcname})
		}
	}
	if err != nil {
		resp.Error(500, err.Error())
		return
	} else {
		resp.Success(tmpres)
		return

	}
}

func AddSshtunnel(ctx *gin.Context) {
	resp := util.NewResult(ctx)
	userID, ok := ctx.Get("userid")
	if !ok {
		global.Logger.Error("addsshtunnel: 没有获取到jwt信息")
		resp.Error(500, "没有获取到jwt信息")
	}
	userinfo := myorm.User{}
	if global.DB.Model(&myorm.User{}).First(&userinfo, userID).Error != nil {
		resp.Error(500, "该用户未在db中查到")
		return
	}

	var postInfo sshTunnel

	if err := ctx.ShouldBind(&postInfo); err != nil {
		// fmt.Println(postInfo)
		global.Logger.Error("addsshtunnel: 获取参数失败")
		resp.Error(500, "获取参数失败")
		return
	}
	tmpres := myorm.Conn{}
	// connID, _ := strconv.Atoi(postInfo.ID)
	if global.DB.Model(&myorm.Conn{}).First(&tmpres, postInfo.ID).Error != nil {
		global.Logger.Error("指定的远程服务不存在")
		resp.Error(500, "指定的远程服务不存在")
		return
	}
	global.Logger.Info(userinfo.Username + "申请关联" + tmpres.Svcname)
	if !userinfo.IsAdmin {
		if err := addWorkflow(userinfo.Username, tmpres.Local, tmpres.Svcname); err != nil {
			global.Logger.Error("增加审批工作流失败： " + err.Error())
			resp.Error(500, err.Error())
			return
		}
	} else {
		if err := starttunnel(0, true, tmpres); err != nil {
			global.Logger.Error(tmpres.Svcname + ": 开启隧道失败: " + err.Error())
			resp.Error(500, err.Error())
			return
		}
	}
	resp.Success(nil)
	return
}

func DelSshtunnel(ctx *gin.Context) {
	resp := util.NewResult(ctx)
	userID, ok := ctx.Get("userid")
	if !ok {
		resp.Error(500, "没有获取到jwt信息")
		return
	}
	var delInfo sshTunnel
	if err := ctx.ShouldBind(&delInfo); err != nil {
		// fmt.Println(err)
		global.Logger.Error(err.Error())
		resp.Error(500, "获取参数失败")
		return
	}
	// fmt.Println("del args: ", delInfo.Remote)
	tmpuser := myorm.User{}
	if err := global.DB.Model(&myorm.User{}).Where("id = ?", userID).Preload("Conn").First(&tmpuser).Error; err != nil {
		global.Logger.Error("db查找用户出错" + err.Error())
		resp.Error(500, "db查找用户出错")
		return
	}
	f := false
	selectConn := myorm.Conn{}
	for _, k := range tmpuser.Conn {
		if delInfo.ID == k.ID {
			selectConn = k
			f = true
		}
	}
	if !f {
		resp.Error(404, "指定的隧道未和你关联授权")
		return
	} else {
		deluser := myorm.User{}
		if err := global.DB.Model(&myorm.User{}).Where("id = ?", userID).First(&deluser).Error; err != nil {
			global.Logger.Error("db查找用户出错: " + err.Error())
			resp.Error(500, "db查找用户出错")
			return
		}
		if err := global.DB.Model(&deluser).Association("Conn").Delete(&selectConn); err != nil {
			global.Logger.Error("取消关联更新db失败: " + err.Error())
			resp.Error(500, "取消关联更新db失败")
			return
		}
		global.Logger.Info(tmpuser.Username + "不关联" + selectConn.Svcname + "成功")
		// 判断 该conn是否还存在其他用户在使用，无则本地端口close，并重置local为“0”
		// tmpconn := myorm.Conn{}
		// if err := global.DB.Model(&myorm.Conn{}).Where(&selectConn).First(&tmpconn).Error; err != nil {
		// 	global.Logger.Error(err.Error())
		// 	resp.Error(500, err.Error())
		// 	return
		// }
		if global.DB.Model(&selectConn).Association("User").Count() == 0 {
			global.Logger.Info(selectConn.Local + " 端口已无人员使用，开始关闭")
			global.Logger.Info(global.GlobalSshtunnelInfo[selectConn.Local])
			// (*(selectConn.St[0])).Close()
			(*(global.GlobalSshtunnelInfo[selectConn.Local])).Close()
			global.Logger.Info(selectConn.Local + " 端口关闭成功")
			if err := global.DB.Model(&myorm.Conn{}).Where("id = ?", delInfo.ID).Update("local", "").Error; err != nil {
				global.Logger.Error("重置local失败: " + err.Error())
				resp.Error(500, "重置local失败")
				return
			}
			global.Logger.Info("重置db中的local为空成功")
		}
		resp.Success(nil)
		return
	}
}
