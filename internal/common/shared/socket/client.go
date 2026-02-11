package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub    Manager
	Conn   *websocket.Conn
	SendCh chan []byte
	UserID string
}

func NewClient(hub Manager, conn *websocket.Conn, userID string) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		SendCh: make(chan []byte, 256),
		UserID: userID,
	}
}

func (c *Client) GetUserID() string { return c.UserID }
func (c *Client) Send(msg []byte)   { c.SendCh <- msg }

// Đọc tin nhắn từ Client gửi lên (nếu có)
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// DÒNG NÀY LÀ HÀM NHẬN CỦA SOCKET THUẦN TÚY
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Server đã nhận được "message" (byte array)
        // Tại đây bạn có thể xử lý:
        // - Nếu là tin nhắn chat -> Gọi ChatUseCase (Nếu làm thuần Socket)
        // - Nếu là lệnh Subscribe -> Gọi Hub.Subscribe (Như mô hình Hybrid)
		// Ở đây có thể handle message từ client gửi lên nếu cần
	}
}

// Gửi tin nhắn từ Server xuống Client
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendCh:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Gửi hết các message đang đợi trong channel cùng 1 lúc (Optimization)
			n := len(c.SendCh)
			for i := 0; i < n; i++ {
				w.Write(<-c.SendCh)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (c *Client) Close() {
	// Đóng kết nối WebSocket
	// Việc này sẽ khiến ReadPump và WritePump bị lỗi (error) và thoát vòng lặp -> giải phóng Goroutine
	c.Conn.Close()
}
