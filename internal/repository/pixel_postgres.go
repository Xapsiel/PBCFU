package repository

import (
	"database/sql"
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/jmoiron/sqlx"
)

type PixelPostgres struct {
	db *sqlx.DB
}

func NewPixelPostgres(db *sqlx.DB) *PixelPostgres {
	return &PixelPostgres{db: db}
}
func (p *PixelPostgres) GetPixels() ([]dewu.Pixel, error) {
	var pixels []dewu.Pixel
	query := "SELECT x, y, id, color FROM pixels"
	err := p.db.Select(&pixels, query)
	if err != nil {
		return nil, err
	}
	return pixels, nil
}
func (p *PixelPostgres) UpdatePixel(pixel dewu.Pixel) error {
	query := fmt.Sprintf(
		"INSERT INTO pixels (x, y, id, color) VALUES (%d, %d, %d, '%s') ON CONFLICT (x, y) DO UPDATE SET id = %d, color = '%s'",
		pixel.X, pixel.Y, pixel.ID, pixel.Color, pixel.ID, pixel.Color)
	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (p *PixelPostgres) GetLastClick(userID int) (int, error) {
	var lastClick int // Переменная для хранения последнего клика

	query := "SELECT lastclick FROM users WHERE id = $1"
	row := p.db.QueryRow(query, userID) // Используем параметризацию запросов для предотвращения SQL-инъекций
	err := row.Scan(&lastClick)         // Сканируем результат в переменную lastClick
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Если пользователь не найден, возвращаем 0
		}
		return 0, err // Возвращаем ошибку, если произошла другая ошибка
	}
	return lastClick, nil // Возвращаем последний клик
}
func (p *PixelPostgres) UpdateClick(userID int, clickValue int) error {
	query := "UPDATE users SET lastclick = $1 WHERE id = $2"
	_, err := p.db.Exec(query, clickValue, userID) // Выполняем запрос на обновление
	return err                                     // Возвращаем ошибку, если произошла ошибка
}
