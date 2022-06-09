package controller

import (
	"fmt"
	"sshtunnelweb/global"
	"sshtunnelweb/myorm"
	"sshtunnelweb/myorm/resps"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

func addWorkflow(username, localport, svcname string) error {
	tmpWorkflow := myorm.Workflow{
		Username:  username,
		Localport: localport,
		Svcname:   svcname,
	}
	err := global.DB.Model(&myorm.Workflow{}).Create(&tmpWorkflow).Error
	return err
}

func changeWorkflow(ID uint) error {
	err := global.DB.Model(&myorm.Workflow{}).Where("id = ?", ID).Updates(&myorm.Workflow{Pass: true}).Error
	return err
}

func starttunnel(workflowID uint, isadmin bool, tmpres myorm.Conn) error {
	tmpconn := tmpres
	tmpworkflow := myorm.Workflow{}
	// 不是admin
	if !isadmin {
		if err := global.DB.Model(&myorm.Workflow{}).Where("id = ?", workflowID).First(&tmpworkflow).Error; err != nil {
			global.Logger.Error(err.Error())
			return err
		}
		tmp := myorm.Conn{}
		if err := global.DB.Model(&myorm.Conn{}).Where(&myorm.Conn{Svcname: tmpworkflow.Svcname}).First(&tmp).Error; err != nil {
			global.Logger.Error(err.Error())
			return err
		}
		tmpconn = tmp
	}
	// 开启连接中转机
	if err := global.StartST(); err != nil {
		global.Logger.Error("连接中转机拨号失败")
		return err
	}
	// 判断要转发的远端内网主机和端口，是否能正常打开
	remote, err := global.ST.Dial("tcp", tmpconn.Remote)
	// fmt.Println("remote addr ", &remote)
	// defer remote.Close()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("远程目标%s不可用 ", tmpconn.Remote))
		global.Logger.Error(err.Error())
		return fmt.Errorf(tmpconn.Remote + "：远端主机不可用")
	}
	global.Logger.Info(tmpconn.Remote + " 可用")
	remote.Close()
	if tmpconn.Local == "" { // 第一次申请，打开端口
		global.Logger.Info(tmpconn.Svcname + ": 第一次开启端口")
		local := util.GetOnePort()
		local_listen, err := global.StartLocalListen(":" + local)
		if err != nil {
			global.Logger.Error("本地监听失败")
			return fmt.Errorf("本地监听失败")
		}
		global.Logger.Info("本地监听 " + local + " 成功")
		global.GlobalSshtunnelInfo[local] = &local_listen
		global.Logger.Info(global.GlobalSshtunnelInfo)
		if global.DB.Model(&myorm.Conn{}).Where("id = ?", tmpconn.ID).Updates(myorm.Conn{Local: local}).Error != nil {
			local_listen.Close()
			global.Logger.Error("db 更新信息识别")
			return fmt.Errorf("db 更新信息识别")
		}
		// 关联 user 和  conn
		if !isadmin {
			tmpuser := myorm.User{}
			if err := global.DB.Model(&myorm.User{}).Where(myorm.User{Username: tmpworkflow.Username}).First(&tmpuser).Error; err != nil {
				local_listen.Close()
				global.Logger.Error(err.Error())
				return fmt.Errorf(err.Error())
			}
			if global.DB.Model(&tmpuser).Association("Conn").Append(&tmpconn) != nil {
				local_listen.Close()
				global.Logger.Error(tmpuser.Username + " 关联 " + tmpconn.Local + " 失败")
				return fmt.Errorf(tmpuser.Username + " 关联 " + tmpconn.Local + " 失败")
			}
			if _, ok := global.LocalPortAndUserIP[local]; ok {
				global.LocalPortAndUserIP[local][tmpuser.Ip] = "1"
			} else {
				global.LocalPortAndUserIP[local] = map[string]string{tmpuser.Ip: "1"}
			}
			global.Logger.Info(global.LocalPortAndUserIP)
			global.Logger.Info(tmpuser.Username + " 关联 " + local + " 成功")
		}
		go global.StartTunnel(local_listen, tmpconn.Remote, global.ST)
	} else { // 已经存在该服务的转发了，直接建立申请人和 已存在conn的连接关系
		// 关联 user 和  conn
		// admin 不用关联
		if !isadmin {
			tmpuser := myorm.User{}
			if err := global.DB.Model(&myorm.User{}).Where(myorm.User{Username: tmpworkflow.Username}).First(&tmpuser).Error; err != nil {
				global.Logger.Error(err.Error())
				return fmt.Errorf(err.Error())
			}
			if global.DB.Model(&tmpuser).Association("Conn").Append(&tmpconn) != nil {
				global.Logger.Error(tmpuser.Username + " 关联 " + tmpconn.Local + " 失败")
				return fmt.Errorf(tmpuser.Username + " 关联 " + tmpconn.Local + " 失败")
			}
			if _, ok := global.LocalPortAndUserIP[tmpconn.Local]; ok {
				global.LocalPortAndUserIP[tmpconn.Local][tmpuser.Ip] = "1"
			} else {
				global.LocalPortAndUserIP[tmpconn.Local] = map[string]string{tmpuser.Ip: "1"}
			}
			global.Logger.Info(global.LocalPortAndUserIP)
			global.Logger.Info(tmpuser.Username + " 关联 " + tmpconn.Local + " 成功")
		}
	}
	return nil
}

func ListWorkflow(c *gin.Context) {
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

	tmpFind := []resps.Workflow{}
	if err := global.DB.Model(&myorm.Workflow{}).Where("pass = ?", false).Find(&tmpFind).Error; err != nil {
		resp.Success(nil)
		return
	}
	resp.Success(&tmpFind)
	return
}

func ChangeOnWorkflow(c *gin.Context) {
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
	type addInfo struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}
	var postInfo addInfo
	if err := c.ShouldBind(&postInfo); err != nil {
		global.Logger.Error(err.Error())
		resp.Error(500, "获取参数失败")
		return
	}
	if err := starttunnel(postInfo.ID, false, myorm.Conn{}); err != nil {
		resp.Error(500, err.Error())
		return
	}
	if err := changeWorkflow(postInfo.ID); err != nil {
		resp.Error(500, err.Error())
		return
	}
	resp.Success(nil)
	return
}
