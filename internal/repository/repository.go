package repository

import (
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/jmoiron/sqlx"
)

// Pixel interface
type Pixel interface {
	GetPixels() ([]dewu.Pixel, error)
	UpdatePixel(pixel dewu.Pixel) error
	GetLastClick(userID int) (int, error) // Добавляем метод
	UpdateClick(userID int, clickValue int) error
}

// User interface
type User interface {
	CreateUser(user dewu.User) (int, error)
	GetUser(login, password string) (dewu.User, error)
	Exist(int, string) (bool, error)
}

type Repository struct {
	Pixel
	User
}

// NewRepository создает новый репозиторий
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:  NewUserPostgres(db),
		Pixel: NewPixelPostgres(db),
	}
}
