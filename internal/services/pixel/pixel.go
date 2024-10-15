package pixel

import (
	"database/sql"
	pixels "dewu/internal/models/pixel"
)

type PixelService struct {
	Pixel *pixels.Pixel
	DB    *sql.DB
}

func New(x, y, owner int, color string, db *sql.DB) *PixelService {
	pixel := pixels.New(x, y, owner, color)
	return &PixelService{Pixel: pixel, DB: db}
}

func (ps *PixelService) Fill() error {
	return ps.Pixel.Fill(ps.DB)
}

func (ps *PixelService) GetPixels() ([]pixels.Pixel, error) {
	return ps.Pixel.GetPixels(ps.DB)
}
