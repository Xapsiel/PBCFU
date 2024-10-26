package service

import (
	"fmt"
	"github.com/Xapsiel/PBCFU/internal/repository"
)

const (
	DefaultUser = 0
	AdminUser   = 1
)

type AdminService struct {
	user repository.User
}

func NewAdminService(repo repository.User) *AdminService {
	return &AdminService{user: repo}
}
func (a *AdminService) IsAdmin(token string) (bool, error) {
	id, login, _, err := (&UserService{}).ParseToken(token)
	if err != nil {
		return false, err
	}
	result, perm, err := a.user.Exist(id, login)
	if err != nil {
		return false, err
	}
	if perm != AdminUser {
		return false, fmt.Errorf("not admin")
	}
	return result, nil
}
