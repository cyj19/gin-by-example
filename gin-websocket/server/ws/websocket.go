/**
 * @Author: cyj19
 * @Date: 2022/5/24 14:15
 */

package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

// ClientManager websocket客户端管理
type ClientManager struct {
	Clients    map[string]*Client // 客户端集合
	BroadCast  chan []byte        // 广播管道
	Register   chan *Client       // 注册管道
	Unregister chan *Client       // 注销管道
}

// Client websocket客户端
type Client struct {
	ID       string
	Conn     *websocket.Conn // websocket连接
	SendChan chan []byte     // 发送管道
}

// Message 消息体
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = &ClientManager{
	Clients:    make(map[string]*Client),
	BroadCast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (m *ClientManager) Close() {
	close(m.Register)
	close(m.Unregister)
	close(m.BroadCast)
	m.Clients = nil
}

// Start 启动websocket服务
func (m *ClientManager) Start(ctx context.Context) {
	defer func() {
		m.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-m.Register:
			log.Println("注册 ", conn.ID)
			m.Clients[conn.ID] = conn
			//msg, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
		case conn := <-m.Unregister:
			if _, ok := m.Clients[conn.ID]; ok {
				log.Println("注销  ", conn.ID)
				// 关闭发送管道
				close(conn.SendChan)
				delete(m.Clients, conn.ID)
			}
		case msg := <-m.BroadCast:
			for _, conn := range m.Clients {
				select {
				case conn.SendChan <- msg:
				default:
					// 无法发送数据，注销该连接
					close(conn.SendChan)
					delete(m.Clients, conn.ID)
				}
			}
		}
	}
}

// Send 向客户端发送数据
func (c *Client) Send() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.SendChan:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}

	}
}

func (c *Client) Receive() {
	defer func() {
		Manager.Unregister <- c
	}()

	for {
		mType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("type:%d msg:%s \n", mType, string(msg))
		// TO DO ...
	}
}
