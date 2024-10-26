package handler

import (
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input dewu.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "success"})
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password"  binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, id, err := h.service.GenerateToken(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "id": id})
}
func (h *Handler) lastClick(c *gin.Context) {
	var LastClick struct {
		Time int `json:"time"`
		ID   int `json:"id"`
	}
	// Получение данных из формы
	if err := c.Bind(&LastClick); err != nil {
		c.String(http.StatusBadRequest, "Неверные данные запроса")
		return
	}
	lastclick, err := h.service.Pixel.GetLastClick(LastClick.ID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "lastclick": lastclick})
}

func (h *Handler) validateToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization") // Получаем токен из заголовка
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"isValid": false})
		return
	}

	// Убираем "Bearer " из токена
	token = token[len("Bearer "):]
	id, login, _, err := h.service.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"isValid": false})
	}
	result, _, err := h.service.Exist(id, login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"isValid": false})
	}
	if !result {
		c.JSON(http.StatusUnauthorized, gin.H{"isValid": false})
	}
	c.JSON(http.StatusOK, gin.H{"isValid": true})
}
