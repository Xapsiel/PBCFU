package models

import (
	"database/sql"
	"dewu/internal/database/postgresql"
)

type Pixel struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Owner int    `json:"owner"`
	Color string `json:"color"`
}

func New(X, Y int, Owner int, color string) *Pixel {

	return &Pixel{X: X, Y: Y, Owner: Owner, Color: color}
}

func (p *Pixel) Fill(db *sql.DB) error {
	repo := postgresql.Repo{DB: db}
	err := repo.Fill(p.X, p.Y, p.Owner, p.Color)
	if err != nil {
		return err
	}
	return nil
}
func (p *Pixel) GetPixels(db *sql.DB) ([]Pixel, error) {
	repo := postgresql.Repo{DB: db}
	pixelsInt, err := repo.GetPixels()
	pixels := make([]Pixel, len(pixelsInt))
	if err != nil {
		return nil, err
	}
	for i := range pixelsInt {
		pixels[i] = *New(pixelsInt[i].X, pixelsInt[i].Y, pixelsInt[i].Owner, pixelsInt[i].Color)
	}
	return pixels, err

}

//func (c *Card) Update(db *sql.DB) error {
//	repo := postgresql.Repo{DB: db}
//	err := repo.UpdateCard(c.Name, c.Price, c.Category)
//	if err != nil {
//		return err
//	}
//	return nil
//}
