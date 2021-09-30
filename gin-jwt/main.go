package main

import (
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

/*
	token中间件gin-jwt的使用
	官方地址：https://github.com/appleboy/gin-jwt
*/

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Username  string
	FirstName string
	LastName  string
}

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(http.StatusOK, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).Username,
		"text":     "Hello World",
	})
}

func main() {
	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8888"
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "gin jwt example",    // jwt标识
		Key:             []byte("secret key"), // 服务端密钥
		Timeout:         time.Hour,            // token过期时间
		MaxRefresh:      time.Hour,            // token最大刷新时间(RefreshToken=Timeout+MaxRefresh)
		IdentityKey:     identityKey,
		PayloadFunc:     payloadFunc,                                        // 添加额外业务相关的信息
		IdentityHandler: identityHandler,                                    // 解析claims
		Authenticator:   authenticator,                                      // 在登录接口中使用的验证方法，并返回验证成功后的用户对象
		Authorizator:    authorizator,                                       // 登录后其他接口验证传入的token方法
		Unauthorized:    unauthorized,                                       // token校验失败处理
		LoginResponse:   loginResponse,                                      // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                     // 登出后的响应
		RefreshResponse: refreshResponse,                                    // 刷新token的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // 依次从请求头、请求参数、cookie中寻找token
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v \n", claims)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "page not found",
		})
	})

	auth := r.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}

}

// 添加额外业务相关的信息
func payloadFunc(data interface{}) jwt.MapClaims {
	// 进行类型断言
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

// 从jwt令牌的声明获取用户标识并传递给授权者
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Username: claims[identityKey].(string),
	}
}

// 登录处理
func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &User{
			Username:  userID,
			FirstName: "vagary",
			LastName:  "yu",
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

// 认证成功处理
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Username == "admin" {
		return true
	}
	return false
}

// 认证失败处理
func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	c.JSON(code, gin.H{
		"token":   token,
		"expires": expires,
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": "logout success",
	})
}

func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	c.JSON(code, gin.H{
		"token":   token,
		"expires": expires,
	})
}
