package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	一切都从hello word说起!
	1. 安装gin:  go get -u github.com/gin-gonic/gin
	2. 创建路由
	3. 添加路由规则和执行函数
	4. 启动服务
*/

const (
	addr = "127.0.0.1:8080"
)

func main() {
	// 创建路由
	g := gin.Default()
	// 绑定路由规则和执行函数
	g.GET("/hello", helloHandler)
	// 启动服务
	if err := http.ListenAndServe(addr, g); err != nil {
		log.Fatalf("服务异常退出：%v \n", err)
	}
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"date": "hello gin",
	})
}
