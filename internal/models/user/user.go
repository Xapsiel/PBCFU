package user

import (
	"crypto/sha1"
	"database/sql"
	"dewu/internal/database/postgresql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	id             int
	Login          string `json:"login" binding:"required"`
	Email          string `json:"email"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password"`
}

const (
	salt          = "tklw12hfoiv3pjihu5u521jofc29urji"
	signingKey    = "gag2rp1jkr21fvi0jio2jqfwcpkkngjy2t0tfp"
	valid_symbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type tokenClaims struct {
	jwt.StandardClaims
	AuthId    int `json:"AuthId"`
	LastClick int `json:"LastClick"`
}

func New(login, password, repeatPassword, email string) *User {
	return &User{Login: login, Password: password, RepeatPassword: repeatPassword, Email: email}
}

func (u *User) SignUp(db *sql.DB) error {
	repo := postgresql.Repo{DB: db}
	u.Password = makeHash(u.Password)

	err := repo.SignUp(u.Login, u.Email, u.Password)
	if err != nil {
		return err
	}

	return nil

}
func makeHash(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
func (u *User) SignIn(db *sql.DB) (string, int, error) {
	repo := postgresql.Repo{DB: db}

	u.Password = makeHash(u.Password)
	id, lastclick, err := repo.SignIn(u.Login, u.Password)
	if err != nil {
		return "", 0, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 + time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		lastclick,
	})
	res, err := token.SignedString([]byte(signingKey))
	return res, id, err
}
func GetLastClick(id int, db *sql.DB) (int, error) {
	repo := postgresql.Repo{DB: db}
	lastClick, err := repo.LastClick(id)
	if err != nil {
		return 0, err
	}
	return lastClick, nil
}
func UpdateLastClick(id, lastclick int, db *sql.DB) error {
	repo := postgresql.Repo{DB: db}
	return repo.UpdateClick(id, lastclick)
}
func ParseToken(accessToken string) (int, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, 0, err
	}
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		id := claims.AuthId
		lastclick := claims.LastClick
		return id, lastclick, nil
	}
	return 0, 0, errors.New("invalid token")

}
func (u *User) Verify() error {
	repo := postgresql.Repo{}
	return repo.VerifyStudent(u.Login, u.Password)
}
