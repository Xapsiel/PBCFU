package repository

import (
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/service/log"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
func (p *UserPostgres) CreateUser(user dewu.User) (int, error) {
	// Проверяем, существует ли пользователь с таким логином или электронной почтой
	var existingUser dewu.User
	query := "SELECT * FROM users WHERE login = $1 OR email = $2"
	err := p.db.Get(&existingUser, query, user.Login, user.Email)

	if err == nil {
		// Если ошибка не произошла, значит пользователь существует
		return 0, fmt.Errorf("Такой пользователь уже существует")
	}

	// Если пользователь не найден, продолжаем создание
	var id int
	insertQuery := "INSERT INTO users (login, email, password) VALUES ($1, $2, $3) returning id"
	row := p.db.QueryRow(insertQuery, user.Login, user.Email, user.Password)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("Ошибка создания пользователя")
	}
	log.Logger.Print(id, "Пользователь был создан")
	return id, nil
}
func (p *UserPostgres) GetUser(login, password string) (dewu.User, error) {
	var user dewu.User
	query := fmt.Sprintf("SELECT * FROM users WHERE login = $1 AND password = $2")
	err := p.db.Get(&user, query, login, password)
	if err != nil {
		return dewu.User{}, fmt.Errorf("Ошибка авторизации")
	}
	log.Logger.Print(user.ID, "Получение информации о пользователе")
	return user, nil
}
func (p *UserPostgres) Exist(id int, login string) (bool, uint, error) {
	var user dewu.User
	query := fmt.Sprintf("SELECT * FROM users WHERE login = $1 AND id = $2")
	err := p.db.Get(&user, query, login, id)
	if err != nil {
		return false, 0, fmt.Errorf("Пользователь не найден")
	}
	log.Logger.Print(user.ID, "Пользователь существует")
	return true, user.Permissions, nil
}
