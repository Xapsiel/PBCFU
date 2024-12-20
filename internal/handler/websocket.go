package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // При необходимости замените проверку на более строгую
	},
}

func (h *Handler) HandleWebSocketConnection(c *gin.Context) {
	// Обновляем соединение до WebSocket
	token := c.Query("token")
	var perm uint = 0
	flag, err := h.service.Admin.IsAdmin(token)
	if err != nil {
		perm = 0
	}
	if flag {
		perm = 1
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to upgrade connection to WebSocket")
		return
	}
	defer conn.Close()

	// Добавляем клиента в WebSocketService
	h.service.AddClient(conn)
	// Обрабатываем соединение
	h.service.HandleConnection(conn, perm)
}
