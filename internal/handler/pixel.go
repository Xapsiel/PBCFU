package handler

import (
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getPixels(c *gin.Context) {
	pixels, err := h.service.Pixel.GetPixels()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if pixels == nil {
		pixels = make([]dewu.Pixel, 0)
	}
	c.JSON(http.StatusOK, pixels)
}

func (h *Handler) print(c *gin.Context) {

}
