package socket

import (
	"log"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RealtimeHandler chịu trách nhiệm upgrade HTTP -> WS
type RealtimeHandler struct {
	manager socket.Manager
}

func NewRealtimeHandler(manager socket.Manager) *RealtimeHandler {
	return &RealtimeHandler{manager: manager}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *RealtimeHandler) HandleWS(c *gin.Context) {
	token := c.Query("token")
	log.Printf("WebSocket connection attempt with token: %s", token)
	// TODO: Validate Token lấy UserID từ Auth Shared Kernel
	userID := "user-demo-123"

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := socket.NewClient(h.manager, conn, userID)
	h.manager.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
