package router

import (
	"net/http"
	"runtime/debug"
	"sshtunnelweb/app/htmlresource"
	ctl "sshtunnelweb/controller"
	"sshtunnelweb/global"
	"sshtunnelweb/middleware"
	"sshtunnelweb/util"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	// r := gin.Default()
	r := gin.New()

	r.Use(middleware.AccessLog(), Recover)
	r.Use(middleware.Cors()) // options通过
	//处理异常
	// r.NoRoute(HandleNotFound)
	r.NoMethod(HandleNotFound)
	r.StaticFS("/assets", http.FS(htmlresource.NewResource())) // 访问/assets的  不需要jwt

	//访问首页/的设置
	index := ctl.NewHtmlHandler()
	r.GET("/", index.Index)
	r.NoRoute(index.RedirectIndex)

	baseGroupRouter := r.Group("/svc")
	{

		baseApi := ctl.Base{}
		baseGroupRouter.POST("/login", baseApi.Login)
		baseGroupRouter.POST("/register", baseApi.Register)

		sshtunnel := baseGroupRouter.Group("/sshtunnel")
		sshtunnel.Use(middleware.JWTAuth())
		{
			sshtunnel.GET("/list", ctl.ListSshtunnel)
			sshtunnel.POST("/add", ctl.AddSshtunnel)
			sshtunnel.POST("/del", ctl.DelSshtunnel)
			remote := sshtunnel.Group("/remote")
			{
				remote.GET("/list", ctl.ListSshtunnelRemote)
				remote.GET("/list/select", ctl.ListSshtunnelRemoteSelect)
				remote.POST("/add", ctl.AddSshtunnelRemote)
				remote.POST("/update", ctl.UpdateSshtunnelRemote)
			}
			local := sshtunnel.Group("/localport")
			{
				local.GET("/list", ctl.ListLocalPort)
				local.POST("/list/user", ctl.ListOneUserLocalPort)
				local.POST("/del/user", ctl.DelOneUserLocalPort)
			}
			workflow := sshtunnel.Group("/workflow")
			{
				workflow.POST("/list", ctl.ListWorkflow)
				workflow.POST("/update", ctl.PassOnWorkflow)
				workflow.POST("/reject", ctl.RejectOnWorkflow)
			}
		}

		sshinfo := baseGroupRouter.Group("/sshinfo")
		sshinfo.Use(middleware.JWTAuth())
		{
			sshinfo.GET("/list", ctl.ListSshinfo)
			sshinfo.POST("/add", ctl.AddSshinfo)
			sshinfo.POST("/update", ctl.UpdateSshinfo)
		}
		opuser := baseGroupRouter.Group("/users")
		opuser.Use(middleware.JWTAuth())
		{
			opuser.GET("/list", ctl.ListUsers)
			opuser.POST("/del", ctl.DelUser)
			opuser.POST("/update", ctl.UpdateUser)
		}
	}

	return r
}

//404
func HandleNotFound(c *gin.Context) {
	global.Logger.Errorf("handle not found: %v", c.Request.RequestURI)
	global.Logger.Errorf("stack: %v", string(debug.Stack()))
	util.NewResult(c).Error(404, "资源未找到")
	return
}

// 500
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			//log.Printf("panic: %v\n", r)
			global.Logger.Errorf("panic: %v", r)
			//log stack
			global.Logger.Errorf("stack: %v", string(debug.Stack()))
			//print stack
			debug.PrintStack()
			//return
			util.NewResult(c).Error(500, "服务器内部错误")
		}
	}()
	//继续后续接口调用
	c.Next()
}
