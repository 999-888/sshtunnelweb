package myjwt

import (
	"errors"
	// "fmt"
	"sshtunnelweb/global"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	Username string
	ID       uint
	IP       string
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	// fmt.Println("port: ", global.CF.Run.Port)
	// fmt.Println("addtime: ", global.CF.Jwt.ExpiresTime, "-")
	ExpiresTime, _ := strconv.Atoi(global.CF.Jwt.ExpiresTime)
	et := time.Now().Unix() + int64(ExpiresTime)
	BufferTime, _ := strconv.Atoi(global.CF.Jwt.BufferTime)
	// fmt.Println("addtime: ", global.CF.Jwt.ExpiresTime, "-", et)
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(BufferTime), // 缓冲时间5m 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, // 签名生效时间
			ExpiresAt: et,
			Issuer:    global.CF.Jwt.Issuer, // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
// func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
// 	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
// 		return j.CreateToken(claims)
// 	})
// 	return v.(string), err
// }

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	// fmt.Println("解析token失败： ", err.Error())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
