package handler

import (
	"github.com/Xapsiel/PBCFU/internal/service/log"
	"github.com/gin-gonic/gin"
)

type error struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Logger.Warn(-1, message)
	c.AbortWithStatusJSON(statusCode, error{message})
}
