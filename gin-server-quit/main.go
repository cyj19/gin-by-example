/**
 * @Author: cyj19
 * @Date: 2021/12/15 16:41
 */

/*
 *  通过例子实现优雅关闭服务
 *  参考地址：https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go
 *  思路：监听系统信号发送到c
 *  signal.Notify函数让signal包将输入信号转发到channel。如果没有列出要传递的信号，会将所有输入信号传递到channel；否则只传递列出的输入信号。
 *	signal包不会为了向channel发送信息而阻塞（就是说如果发送时channel阻塞了，signal包会直接放弃）
 *  调用者应该保证channel有足够的缓存空间可以跟上期望的信号频率。对使用单一信号用于通知的通道，缓存为1就足够了
 */

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	g := gin.Default()
	g.GET("/hello", HelloHandle)

	// 启动服务
	srv := &http.Server{
		Addr:    ":8888",
		Handler: g,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %#v\n", err)
		}
	}()

	// 监听系统打断信号
	quit := make(chan os.Signal, 1)
	// SIGINT: ctrl+C触发  SIGTERM: 结束程序触发
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 保留5秒处理业务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %#v\n", err)
	}

	log.Println("Server is shutdown")
}

func HelloHandle(c *gin.Context) {
	// 模拟需要耗时长的操作
	time.Sleep(4 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "Elegant shutdown server",
	})
}
