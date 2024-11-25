package repository

import (
	"database/sql"
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/service/log"
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
		return nil, fmt.Errorf("Ошибка получения пикселей")
	}
	log.Logger.Print(-1, "Получение данных всех пикселей")
	return pixels, nil
}
func (p *PixelPostgres) UpdatePixel(pixel dewu.Pixel) error {
	query := fmt.Sprintf(
		"INSERT INTO pixels (x, y, id, color) VALUES (%d, %d, %d, '%s') ON CONFLICT (x, y) DO UPDATE SET id = %d, color = '%s'",
		pixel.X, pixel.Y, pixel.ID, pixel.Color, pixel.ID, pixel.Color)
	_, err := p.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Ошибка обновления данных о пикселе")
	}
	log.Logger.Print(pixel.ID, fmt.Sprintf("Обновление данных пикселя [%d;%d] на следующий цвет [%s]", pixel.X, pixel.Y, pixel.Color))
	return nil
}
func (p *PixelPostgres) GetLastClick(userID int) (int, error) {
	var lastClick int // Переменная для хранения последнего клика

	query := "SELECT lastclick FROM users WHERE id = $1"
	row := p.db.QueryRow(query, userID) // Используем параметризацию запросов для предотвращения SQL-инъекций
	err := row.Scan(&lastClick)         // Сканируем результат в переменную lastClick
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("Пользователь не найден") // Если пользователь не найден, возвращаем 0
		}
		return 0, fmt.Errorf("Произошла ошибка получения данных о последних кликах") // Возвращаем ошибку, если произошла другая ошибка
	}
	log.Logger.Print(userID, "Получение таймкода последнего действия")
	return lastClick, nil // Возвращаем последний клик
}
func (p *PixelPostgres) UpdateClick(userID int, clickValue int) error {
	query := "UPDATE users SET lastclick = $1 WHERE id = $2"
	_, err := p.db.Exec(query, clickValue, userID) // Выполняем запрос на обновление
	if err != nil {
		return fmt.Errorf("Ошибка обновления данных об изменении пикселя")
	}
	log.Logger.Print(userID, "Обновление таймкода последнего действия")
	return nil // Возвращаем ошибку, если произошла ошибка
}
