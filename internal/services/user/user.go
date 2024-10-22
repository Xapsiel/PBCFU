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
func (us *UserService) SignIn() (string, int, error) {
	u := user.New(us.Student.Login, us.Student.Password, "", "")
	return u.SignIn(us.DB)
}

func (us *UserService) Verify() error {
	u := user.New(us.Student.Login, us.Student.Password, us.Student.Repeatpassword, us.Student.Email)
	return u.Verify(us.DB)
}

func ParseToken(accessToken string) (string, int, int, error) {
	return user.ParseToken(accessToken)
}
func GetLastClick(id int, db *sql.DB) (int, error) {
	return user.GetLastClick(id, db)
}
func UpdateLastClick(id, lastclick int, db *sql.DB) error {
	return user.UpdateLastClick(id, lastclick, db)
}

func Exists(id int, login string, db *sql.DB) bool {
	u := user.New(login, "", "", "")
	return u.Exist(id, db)
}
