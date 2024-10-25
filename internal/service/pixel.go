package service

import (
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/repository"
)

type PixelService struct {
	repo repository.Pixel
}

func NewPixelService(repo repository.Pixel) *PixelService {
	return &PixelService{repo: repo}
}

func (s *PixelService) GetPixels() ([]dewu.Pixel, error) {
	return s.repo.GetPixels()
}

func (s *PixelService) UpdatePixel(pixel dewu.Pixel) error {
	return s.repo.UpdatePixel(pixel)
}

func (s *PixelService) GetLastClick(userID int) (int, error) {
	return s.repo.GetLastClick(userID)
}
func (s *PixelService) UpdateClick(userID int, clickValue int) error {
	return s.repo.UpdateClick(userID, clickValue)

}
