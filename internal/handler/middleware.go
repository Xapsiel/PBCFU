package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No Authorization header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header")
		return
	}
	userid, login, lastclick, err := h.service.User.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header")
	}
	c.Set("token", headerParts[1])
	c.Set("user", userid)
	c.Set("login", login)
	c.Set("lastclick", lastclick)

}
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
func (h *Handler) isAdmin(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		c.Set("permission", 0)
	}
	tokenstr := token.(string)
	flag, err := h.service.IsAdmin(tokenstr)
	if err != nil {
		c.Set("permission", 0)
	}
	if !flag {
		c.Set("permission", 0)
	}
	c.Set("permission", 1)
	c.Next()
}
