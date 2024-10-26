package repository

import (
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
func (p *UserPostgres) CreateUser(user dewu.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (login, email, password) VALUES ($1, $2, $3) returning id")
	row := p.db.QueryRow(query, user.Login, user.Email, user.Password)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
func (p *UserPostgres) GetUser(login, password string) (dewu.User, error) {
	var user dewu.User
	query := fmt.Sprintf("SELECT * FROM users WHERE login = $1 AND password = $2")
	err := p.db.Get(&user, query, login, password)
	return user, err
}
func (p *UserPostgres) Exist(id int, login string) (bool, uint, error) {
	var user dewu.User
	query := fmt.Sprintf("SELECT * FROM users WHERE login = $1 AND id = $2")
	err := p.db.Get(&user, query, login, id)
	if err != nil {
		return false, 0, err
	}
	return true, user.Permissions, nil
}
