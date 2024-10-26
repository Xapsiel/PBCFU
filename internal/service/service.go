package service

import (
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/repository"
	"github.com/gorilla/websocket"
)

// Pixel - интерфейс для работы с пикселями
type Pixel interface {
	GetPixels() ([]dewu.Pixel, error) // Получение всех пикселей
	GetLastClick(userID int) (int, error)
	UpdatePixel(pixel dewu.Pixel) error
	UpdateClick(userID int, clickValue int) error
}

// User - интерфейс для работы с пользователями
type User interface {
	CreateUser(user dewu.User) (int, error)                           // Создание нового пользователя
	GenerateToken(login string, password string) (string, int, error) // Генерация токена
	ParseToken(token string) (int, string, int, error)                // Парсинг токена
	Exist(int, string) (bool, uint, error)
}

// Websocket - интерфейс для работы с WebSocket
type Websocket interface {
	AddClient(conn *websocket.Conn)                   // Добавление клиента
	HandleConnection(conn *websocket.Conn, perm uint) // Обработка соединения
	BroadcastMessage()                                // Рассылка сообщений всем клиентам
}
type Admin interface {
	IsAdmin(string) (bool, error)
}

// Service объединяет все зависимости приложения
type Service struct {
	Pixel
	User
	Websocket
	Admin
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:      NewUserService(repo.User),
		Pixel:     NewPixelService(repo.Pixel),
		Websocket: NewWebSocketService(repo.Pixel),
		Admin:     NewAdminService(repo.User),
	}
}
