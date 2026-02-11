package socket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9" // Import Redis
)

const (
	PubSubChannel = "social_network_realtime"
)

type Hub struct {
	clients    map[string]map[ClientOps]bool
	register   chan ClientOps
	unregister chan ClientOps
	broadcast  chan Message

	rdb *redis.Client // Redis Connection
	mu  sync.RWMutex
}

// Cần inject Redis Client vào Hub
func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[string]map[ClientOps]bool),
		register:   make(chan ClientOps),
		unregister: make(chan ClientOps),
		broadcast:  make(chan Message),
		rdb:        rdb,
	}
}

func (h *Hub) Run(ctx context.Context) {
	// 1. Subscribe Redis để nhận tin từ các Server khác
	pubsub := h.rdb.Subscribe(ctx, PubSubChannel)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err == nil {
				h.broadcastLocal(message) // Gửi cho user nếu họ đang kết nối tới server này
			}
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			uid := client.GetUserID()
			if _, ok := h.clients[uid]; !ok {
				h.clients[uid] = make(map[ClientOps]bool)
			}
			h.clients[uid][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			uid := client.GetUserID()
			if _, ok := h.clients[uid]; ok {
				delete(h.clients[uid], client)
				client.Close() // Đóng connection an toàn
				if len(h.clients[uid]) == 0 {
					delete(h.clients, uid)
				}
			}
			h.mu.Unlock()

		case <-ctx.Done(): // Graceful Shutdown
			h.mu.Lock()
			for _, clients := range h.clients {
				for client := range clients {
					client.Close()
				}
			}
			h.mu.Unlock()
			return
		}
	}
}

// SendToUser: Thay vì gửi trực tiếp, ta Publish lên Redis
// Để TẤT CẢ server đều biết và gửi cho user (bất kể user đang kết nối server nào)
func (h *Hub) SendToUser(userID string, msg Message) error {
	msg.ToUserID = userID
	data, _ := json.Marshal(msg)
	return h.rdb.Publish(context.Background(), PubSubChannel, data).Err()
}

// Hàm nội bộ: Chỉ gửi cho connection đang nằm trên server này
func (h *Hub) broadcastLocal(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[msg.ToUserID]; ok {
		data, _ := json.Marshal(msg)
		for client := range clients {
			client.Send(data)
		}
	}
}

func (h *Hub) Broadcast(msg Message)  { /* Tương tự SendToUser nhưng không check UserID */ }
func (h *Hub) Register(c ClientOps)   { h.register <- c }
func (h *Hub) Unregister(c ClientOps) { h.unregister <- c }
