package rest

import (
	"database/sql"
	"dewu/internal/config"
	"dewu/internal/services/painter"
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
var broadcast = make(chan pixelClick)

type pixelClick struct {
	*PixSer.PixelService
	lastclick int
}

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
	token, id, err := u.SignIn()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Возврат результата
	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token, "id": id})
}

func handleMessages() {
	for {
		pixel := <-broadcast

		// Логирование данных пикселя

		// Отправка данных пикселя всем клиентам
		for client := range clients {
			if err := client.WriteJSON(pixel); err != nil {
				log.Printf("Ошибка отправки данных клиенту: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func handleConnections(c *gin.Context, db *sql.DB) {
	// Проверяем токен перед подключением к WebSocket
	_, id, _, err := JWTMiddlewareWebSocket(c, db)
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
			X         int    `json:"x"`
			Y         int    `json:"y"`
			Color     string `json:"color"`
			Owner     int    `json:"owner"`
			Lastclick int    `json:"lastclick"`
		}

		// Чтение данных пикселя из вебсокета
		if err := ws.ReadJSON(&pixel); err != nil {
			log.Printf("Ошибка чтения JSON: %v", err)
			break // Break the loop on read error
		}
		fmt.Println(pixel)
		// Создание нового объекта пикселя и проверка пользователя
		pixelSer := *PixSer.New(pixel.X, pixel.Y, pixel.Owner, pixel.Color, db)
		pixelclk := pixelClick{&pixelSer, pixel.Lastclick}
		user.UpdateLastClick(id, pixelclk.lastclick, db)

		// Заполнение пикселя
		if err := pixelSer.Fill(); err != nil {
			log.Printf("Ошибка заполнения пикселя: %v", err)
			continue
		}

		// Отправка пикселя в канал
		broadcast <- pixelclk
	}

	// Удаляем клиента после выхода из цикла
	delete(clients, ws)
}

func Print(c *gin.Context, db *sql.DB) {
	// Определяем структуру для получения пикселей
	var pixels struct {
		ID   int `json:"id"`
		Data []struct {
			X     int    `json:"x"`
			Y     int    `json:"y"`
			Color string `json:"color"`
		} `json:"data"` // Добавляем тег json
	}
	//var d interface{}
	// Пробуем привязать данные из запроса к структуре pixels
	if err := c.BindJSON(&pixels); err != nil {
		c.String(http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	// Извлекаем ID и данные
	id := pixels.ID
	data := pixels.Data

	// Заполняем пиксели в базе данных
	for _, pixel := range data {
		pixelSer := *PixSer.New(pixel.X, pixel.Y, id, pixel.Color, db)
		if err := pixelSer.Fill(); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка заполнения пикселя: %v", err))
			return
		}
	}

	// Возвращаем успешный ответ
	all := painter.New("smile.png")
	for _, row := range all {
		for _, pixel := range row {
			pixelSer := *PixSer.New(pixel.X, pixel.Y, id, pixel.Color, db)
			if err := pixelSer.Fill(); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка заполнения пикселя: %v", err))
				return
			}
		}
	}
	c.String(http.StatusOK, "Пиксели успешно добавлены")

}

func getLastClickHandler(c *gin.Context, db *sql.DB) {
	var LastClick struct {
		Time int `json:"time"`
		ID   int `json:"id"`
	}
	// Получение данных из формы
	if err := c.Bind(&LastClick); err != nil {
		c.String(http.StatusBadRequest, "Неверные данные запроса")
		return
	}
	lastclick, err := user.GetLastClick(LastClick.ID, db)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "lastclick": lastclick})

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
func JWTMiddlewareWebSocket(c *gin.Context, db *sql.DB) (string, int, int, error) {
	token := c.Query("token") // Получаем токен из query-параметров
	if token == "" {
		return "", 0, 0, fmt.Errorf("не передан JWT токен")
	}

	// Здесь разбираем токен, если нужно, его декодируем или парсим
	id, lastclick, err := user.ParseToken(token)
	if err != nil {
		return "", 0, 0, fmt.Errorf("неверный JWT токен")
	}
	return token, id, lastclick, nil
}

// Обработчик соединений по WebSocket

// Запуск сервера
func StartServer(cfg config.Config, db *sql.DB) error {
	router := gin.Default()
	// Использование CORS
	router.Use(corsMiddleware)
	router.POST("/print", func(c *gin.Context) {
		Print(c, db)
	})
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

	router.POST("/getLastClick", func(c *gin.Context) {
		getLastClickHandler(c, db)
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
