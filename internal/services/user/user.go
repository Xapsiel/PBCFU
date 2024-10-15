package user

import (
	"database/sql"
	"dewu/internal/models/user"
)

type UserService struct {
	Student struct {
		Login          string
		Password       string
		Repeatpassword string
		Email          string
	}
	DB *sql.DB
}

func New(login, password, repeatPassword, email string, db *sql.DB) *UserService {
	return &UserService{Student: struct {
		Login          string
		Password       string
		Repeatpassword string
		Email          string
	}{Login: login, Password: password, Repeatpassword: repeatPassword, Email: email}, DB: db}
}

func (us *UserService) SignUp() error {
	u := user.New(us.Student.Login, us.Student.Password, us.Student.Repeatpassword, us.Student.Email)
	err := u.SignUp(us.DB)
	if err != nil {
		return err
	}
	return nil
}
func (us *UserService) SignIn() (string, error) {
	u := user.New(us.Student.Login, us.Student.Password, "", "")
	return u.SignIn(us.DB)
}

func (us *UserService) Verify() error {
	u := user.New(us.Student.Login, us.Student.Password, us.Student.Repeatpassword, us.Student.Email)
	return u.Verify()
}

func ParseToken(accessToken string) (int, error) {
	return user.ParseToken(accessToken)
}
