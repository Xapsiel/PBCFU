package handler

import (
	"github.com/Xapsiel/PBCFU/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(corsMiddleware)
	auth := router.Group("auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
		auth.GET("validateToken", h.validateToken)
	}
	api := router.Group("api", h.userIdentity)
	{
		api.POST("/getLastClick", h.lastClick)
		api.POST("/print", h.print)
	}
	pixels := router.Group("pixels")
	{
		pixels.GET("/getPixels", h.getPixels)
	}

	webhook := router.Group("/webhook")
	{
		webhook.GET("/ws", h.HandleWebSocketConnection)
	}
	return router

}
