package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"net/mail"
	"strings"
	"time"
)

const (
	signingKey = ("afgkasogdnasgvuio2r1jioqwdf89zsfiolkasf")
	salt       = ";knmmm3rjoq; 2vr541jdhaDCGV1UE9PED"
	en_lower   = "abcdefghijklmnopqrstuvwxyz"
	en_upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits     = "0123456789"
	symbols    = "@_"
)

type UserService struct {
	repo repository.User
}
type tokenClaims struct {
	jwt.StandardClaims
	UserId    int    `json:"UserId"`
	Login     string `json:"Login"`
	LastClick int    `json:"LastClick"`
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}
func (u *UserService) CreateUser(user dewu.User) (int, error) {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return 0, fmt.Errorf("Неправильный формат электронной почты")
	}
	err = validateLogin(user.Login)
	if err != nil {
		return 0, err
	}
	if user.Password != user.RepeatPassword {
		return 0, fmt.Errorf("Пароли должны совпадать")
	}
	if _, err := validatePassword(user.Password); err != nil {
		return 0, err
	}
	user.Password = generatePasswordHash(user.Password)
	user.RepeatPassword = user.Password
	return u.repo.CreateUser(user)
}
func (u *UserService) Exist(id int, login string) (bool, uint, error) {
	return u.repo.Exist(id, login)
}
func (u *UserService) GenerateToken(login, password string) (string, int, error) {

	user, err := u.repo.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", 0, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:    user.ID,
		Login:     login,
		LastClick: user.LastClick,
	})
	tokenStr, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", 0, fmt.Errorf("Ошибка генерации токена")
	}
	return tokenStr, user.ID, nil
}
func (u *UserService) ParseToken(accessToken string) (int, string, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, "", 0, errors.New("invalid token")
	}
	return claims.UserId, claims.Login, claims.LastClick, nil
}
func validatePassword(password string) (bool, error) {
	for _, elem := range password {
		if !strings.Contains(en_lower+en_upper+digits+symbols, string(elem)) {
			return false, fmt.Errorf("Пароль может содержать только символы английского алфавита, цифры, @_")
		}
	}
	if len(password) < 8 {
		return false, fmt.Errorf("Пароль должен быть длинной 8 и больше")
	}
	return true, nil
}
func validateLogin(login string) error {
	en_letter_contain := false
	for _, elem := range login {
		if strings.Contains(en_lower+en_upper, string(elem)) {
			en_letter_contain = true
			continue
		} else if strings.Contains(digits, string(elem)) {
			continue
		}
		return fmt.Errorf("Неправильный формат логина пользователя")
	}
	if len(login) < 8 {
		return fmt.Errorf("Минимальная длина логина - 8 букв")
	}
	if !(en_letter_contain) {
		return fmt.Errorf("Неправильный формат логина пользователя")
	}
	return nil
}
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
