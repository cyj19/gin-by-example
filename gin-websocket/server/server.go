/**
 * @Author: cyj19
 * @Date: 2022/5/24 14:14
 */

package main

import (
	"context"
	"github.com/cyj19/gin-by-example/gin-websocket/server/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g := gin.Default()
	g.GET("/ws", wsHandle)
	go ws.Manager.Start(ctx)

	s := &http.Server{
		Addr:           ":8081",
		Handler:        g,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe error: %#v \n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	log.Println("server exit...")
}

func wsHandle(c *gin.Context) {
	wu := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	// 升级为websocket协议
	conn, err := wu.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	client := &ws.Client{ID: xid.New().String(), Conn: conn, SendChan: make(chan []byte)}

	// 注册
	ws.Manager.Register <- client

	go client.Send()
	go client.Receive()
}
