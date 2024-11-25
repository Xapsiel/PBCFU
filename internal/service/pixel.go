package service

import (
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/repository"
)

type PixelService struct {
	repo repository.Pixel
}

const (
	CanvasHeight = 100
	CanvasWidth  = 100
)

func NewPixelService(repo repository.Pixel) *PixelService {
	return &PixelService{repo: repo}
}

func (s *PixelService) GetPixels() ([]dewu.Pixel, error) {
	return s.repo.GetPixels()
}

func (s *PixelService) UpdatePixel(pixel dewu.Pixel) error {
	if (CanvasWidth-pixel.X <= 0) || (CanvasHeight-pixel.Y <= 0) {
		return fmt.Errorf("Ошибка обновления данных о пикселе")
	}
	return s.repo.UpdatePixel(pixel)
}

func (s *PixelService) GetLastClick(userID int) (int, error) {
	return s.repo.GetLastClick(userID)
}
func (s *PixelService) UpdateClick(userID int, clickValue int) error {
	return s.repo.UpdateClick(userID, clickValue)

}
