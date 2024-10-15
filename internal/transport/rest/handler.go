package rest

import (
	"database/sql"
	"dewu/internal/config"
	PixSer "dewu/internal/services/pixel"
	"dewu/internal/services/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Настройка апгрейдера вебсокета
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Разрешаем все запросы
}

// Список клиентов вебсокета и канал для передачи пикселей
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan PixSer.PixelService)

// Обработчик регистрации
func SignUpHandler(c *gin.Context, db *sql.DB) {
	u := user.UserService{DB: db}
	var LOGIN struct {
		Login          string `json:"login"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		RepeatPassword string `json:"repeatpassword"`
	}
	// Получение данных из формы
	if err := c.Bind(&LOGIN); err != nil {
		c.String(http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Проверка совпадения паролей
	if LOGIN.Password != LOGIN.RepeatPassword {
		c.String(http.StatusBadRequest, "Пароли не совпадают")
		return
	}

	// Заполнение структуры пользователя
	u.Student.Login = LOGIN.Login
	u.Student.Email = LOGIN.Email
	u.Student.Password = LOGIN.Password

	// Попытка регистрации
	if err := u.SignUp(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Обработчик входа
func SignInHandler(c *gin.Context, db *sql.DB) {
	u := user.UserService{DB: db}
	var LOGIN struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// Получение данных из формы
	if err := c.Bind(&LOGIN); err != nil {
		c.String(http.StatusBadRequest, "Неверные данные запроса")
		return
	}
	// Получение данных из формы

	// Заполнение структуры пользователя
	u.Student.Password = LOGIN.Password
	u.Student.Login = LOGIN.Login

	// Попытка входа
	token, err := u.SignIn()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Возврат результата
	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func handleMessages() {
	for {
		pixel := <-broadcast

		// Логирование данных пикселя
		fmt.Println(pixel.Pixel.X, pixel.Pixel.Y, pixel.Pixel.Owner, pixel.Pixel.Color)

		// Отправка данных пикселя всем клиентам
		for client := range clients {
			if err := client.WriteJSON(pixel.Pixel); err != nil {
				log.Printf("Ошибка отправки данных клиенту: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// Обработчик получения всех пикселей
func getPixelsHandler(c *gin.Context, db *sql.DB) {
	pixelService := PixSer.PixelService{DB: db}
	pixels, err := pixelService.GetPixels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pixels)
}

// Middleware для CORS
func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Добавляем Authorization

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent) // Обрабатываем предзапрос
		return
	}

	c.Next()
}

// Middleware для JWT, специфичный для WebSocket
func JWTMiddlewareWebSocket(c *gin.Context) (string, error) {
	token := c.Query("token") // Получаем токен из query-параметров
	if token == "" {
		return "", fmt.Errorf("не передан JWT токен")
	}

	// Здесь разбираем токен, если нужно, его декодируем или парсим
	if _, err := user.ParseToken(token); err != nil {
		return "", fmt.Errorf("неверный JWT токен")
	}
	return token, nil
}

// Обработчик соединений по WebSocket
func handleConnections(c *gin.Context, db *sql.DB) {
	// Проверяем токен перед подключением к WebSocket
	_, err := JWTMiddlewareWebSocket(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ошибка обновления до вебсокета:", err)
		return
	}
	defer ws.Close()

	// Добавляем нового клиента
	clients[ws] = true
	defer delete(clients, ws)

	for {
		// Структура для данных пикселя
		var pixel struct {
			X     int    `json:"x"`
			Y     int    `json:"y"`
			Color string `json:"color"`
			Owner int    `json:"owner"`
		}

		// Чтение данных пикселя из вебсокета
		if err := ws.ReadJSON(&pixel); err != nil {
			log.Printf("Ошибка чтения JSON: %v", err)
			continue
		}

		// Создание нового объекта пикселя и проверка пользователя
		pixelSer := *PixSer.New(pixel.X, pixel.Y, pixel.Owner, pixel.Color, db)

		// Заполнение пикселя
		if err := pixelSer.Fill(); err != nil {
			log.Printf("Ошибка заполнения пикселя: %v", err)
			continue
		}

		// Отправка пикселя в канал
		broadcast <- pixelSer
	}
}

// Запуск сервера
func StartServer(cfg config.Config, db *sql.DB) error {
	router := gin.Default()

	// Использование CORS
	router.Use(corsMiddleware)

	// Группа маршрутов для API без JWT аутентификации
	router.POST("/SignUp", func(c *gin.Context) {
		SignUpHandler(c, db)
	})
	router.POST("/SignIn", func(c *gin.Context) {
		SignInHandler(c, db)
	})
	router.GET("/getPixels", func(c *gin.Context) {
		getPixelsHandler(c, db)
	})
	// Группа маршрутов для API с JWT аутентификацией
	webhook := router.Group("/webhook")
	// Маршруты для WebSocket
	webhook.GET("/ws", func(c *gin.Context) {
		handleConnections(c, db)
	})

	// Запуск обработчика сообщений в отдельной горутине
	go handleMessages()

	log.Println("Сервер запущен на порту :8080")
	return router.Run(":8080")

}
