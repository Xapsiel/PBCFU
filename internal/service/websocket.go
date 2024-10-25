package service

import (
	"github.com/Xapsiel/PBCFU/internal/repository"
	"log"

	dewu "github.com/Xapsiel/PBCFU"
	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	pixelService *PixelService

	clients   map[*websocket.Conn]bool
	broadcast chan dewu.PixelClick // Канал для вещания пикселей
}

// Конструктор для WebSocketService
func NewWebSocketService(repo repository.Pixel) *WebSocketService {
	service := &WebSocketService{
		pixelService: NewPixelService(repo),
		clients:      make(map[*websocket.Conn]bool),
		broadcast:    make(chan dewu.PixelClick),
	}

	// Запуск обработки вещания в отдельной горутине
	go service.BroadcastMessage()

	return service
}

// Метод для обработки входящих сообщений WebSocket
func (ws *WebSocketService) HandleConnection(conn *websocket.Conn) {
	ws.AddClient(conn) // Добавляем клиента
	for {
		var pixel dewu.PixelClick

		// Чтение данных пикселя от клиента
		err := conn.ReadJSON(&pixel)

		if err != nil {
			log.Printf("Ошибка чтения JSON: %v", err)
			ws.RemoveClient(conn) // Удаляем клиента в случае ошибки
			break
		}

		ws.pixelService.UpdatePixel(dewu.Pixel{ID: pixel.ID, X: pixel.X, Y: pixel.Y, Color: pixel.Color})
		ws.pixelService.UpdateClick(pixel.ID, pixel.Lastclick)
		// Добавляем пиксель в канал вещания
		ws.broadcast <- pixel
	}
}

// Метод для добавления клиентов в WebSocketService
func (ws *WebSocketService) AddClient(conn *websocket.Conn) {
	ws.clients[conn] = true
}

// Метод для удаления клиента
func (ws *WebSocketService) RemoveClient(conn *websocket.Conn) {
	if _, ok := ws.clients[conn]; ok {
		delete(ws.clients, conn)
		conn.Close()
	}
}

// Обработка вещания и рассылка сообщений всем клиентам
func (ws *WebSocketService) BroadcastMessage() {
	for {
		// Чтение из канала broadcast
		pixel := <-ws.broadcast

		// Рассылка пикселя всем клиентам
		for client := range ws.clients {
			err := client.WriteJSON(pixel)
			if err != nil {
				log.Printf("Ошибка отправки пикселя клиенту: %v", err)
				ws.RemoveClient(client) // Удаляем клиента в случае ошибки отправки
			}
		}
	}
}
