package middleware

import (
	"sshtunnelweb/app/myjwt"
	"sshtunnelweb/global"
	"sshtunnelweb/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := util.NewResult(c)
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get(global.JwtHeaderName)
		if token == "" {
			resp.Error(401, "没有登录")
			c.Abort()
			return
		}
		j := myjwt.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		// fmt.Println(claims, err)
		if err != nil {
			if err == myjwt.TokenExpired {
				resp.Error(403, "授权已过期")
				c.Abort()
				return
			}
			resp.Error(401, err.Error())
			c.Abort()
			return
		}
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开
		//if err, _ = userService.FindUserByUuid(claims.UUID.String()); err != nil {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}

		// 刷新token问题处理  后端主动生成新token，放到本次响应header中，前端监听这个header，有就替换token
		// fmt.Println(claims.ExpiresAt, time.Now().Unix(), claims.BufferTime)
		ExpiresTime, _ := strconv.Atoi(global.CF.Jwt.ExpiresTime)
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + int64(ExpiresTime)
			// claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(global.CF.Jwt.ExpiresTime))).Unix()
			newToken, _ := j.CreateToken(*claims)
			// newClaims, _ := j.ParseToken(newToken)
			c.Header(global.NewJwtHeaderName, newToken)
			// c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
		}
		// c.Set("claims", claims)
		c.Set("userid", claims.ID)
		c.Next()
	}
}
