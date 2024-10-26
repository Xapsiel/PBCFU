package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	signingKey = ("afgkasogdnasgvuio2r1jioqwdf89zsfiolkasf")
	salt       = ";knmmm3rjoq; 2vr541jdhaDCGV1UE9PED"
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
	if user.Password != user.RepeatPassword {
		return 0, fmt.Errorf("password is different")
	}
	user.Password = generatePasswordHash(user.Password)
	user.RepeatPassword = user.Password
	return u.repo.CreateUser(user)
}
func (u *UserService) Exist(id int, login string) (bool, uint, error) {
	return u.repo.Exist(id, login)
}
func (u *UserService) GenerateToken(login string, password string) (string, int, error) {
	//get user from db
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
	return tokenStr, user.ID, err
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

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
