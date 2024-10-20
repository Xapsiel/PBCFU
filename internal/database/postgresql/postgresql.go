package postgresql

import (
	"database/sql"
	"dewu/internal/config"
	"fmt"
	_ "github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}
type RepoPoint struct {
	X     int
	Y     int
	Owner int
	Color string
}

func New(cfg config.DatabaseConfig) (*Repo, error) {
	repo := Repo{}
	var err error
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.DBName)
	repo.SetupDB(cfg)
	repo.DB, err = sql.Open("postgres", connstr)
	return &repo, err
}

func (r *Repo) SetupDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	connstr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", cfg.User, cfg.Password)
	err := r.createDB(cfg.DBName, connstr)
	if err != nil {
		return nil, err
	}
	connstr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.DBName)
	query := "CREATE TABLE IF NOT EXISTS Users(id SERIAL PRIMARY KEY,login text,email text,password text,lastclick bigint DEFAULT 0) "
	err = r.createTable(connstr, query)
	if err != nil {
		return nil, err
	}
	query = `
    CREATE TABLE IF NOT EXISTS Pixels (
        x INT,
        y INT,
        owner INT,
        color VARCHAR(7) NULL,
        PRIMARY KEY (x, y)
    )
`
	err = r.createTable(connstr, query)
	if err != nil {
		return nil, err
	}
	return r.DB, nil
}

func (r *Repo) createDB(dbName, connstr string) error {
	var err error
	r.DB, err = sql.Open("postgres", connstr)
	if err != nil {
		return err
	}
	dublicate := fmt.Sprintf("pq: database \"%s\" already exists", dbName)
	_, err = r.DB.Exec("create database " + dbName)
	if err != nil && err.Error() != dublicate {
		return err
	}
	return nil
}

func (r *Repo) createTable(connstr, query string) error {
	var err error

	r.DB, err = sql.Open("postgres", connstr)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) Fill(x, y, owner int, color string) error {
	query := fmt.Sprintf(
		"INSERT INTO pixels (x, y, owner, color) VALUES (%d, %d, %d, '%s') ON CONFLICT (x, y) DO UPDATE SET owner = %d, color = '%s'",
		x, y, owner, color, owner, color,
	)
	_, err := r.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) LastClick(id int) (int, error) {
	var lastClick int // Переменная для хранения последнего клика

	query := "SELECT lastclick FROM users WHERE id = $1"
	row := r.DB.QueryRow(query, id) // Используем параметризацию запросов для предотвращения SQL-инъекций
	err := row.Scan(&lastClick)     // Сканируем результат в переменную lastClick
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Если пользователь не найден, возвращаем 0
		}
		return 0, err // Возвращаем ошибку, если произошла другая ошибка
	}
	return lastClick, nil // Возвращаем последний клик
}
func (r *Repo) UpdateClick(userID int, clickValue int) error {
	query := "UPDATE users SET lastclick = $1 WHERE id = $2"
	_, err := r.DB.Exec(query, clickValue, userID) // Выполняем запрос на обновление
	return err                                     // Возвращаем ошибку, если произошла ошибка
}
func (r *Repo) GetPixels() ([]RepoPoint, error) {
	var (
		x     int
		y     int
		owner int
		color string
	)
	query := "SELECT x, y, owner, color FROM pixels"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var array []RepoPoint = make([]RepoPoint, 0)
	for rows.Next() {
		err := rows.Scan(&x, &y, &owner, &color)
		if err != nil {
			return nil, err
		}
		array = append(array, RepoPoint{
			X:     x,
			Y:     y,
			Owner: owner,
			Color: color,
		})
	}
	return array, nil
}

func (r *Repo) SignUp(login, email, password string) error {
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 OR email = $2)", login, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if exists {
		return fmt.Errorf("user with login '%s' or email '%s' already exists", login, email)
	}

	// Insert new user if not exists
	query := "INSERT INTO users (login, email, password) VALUES ($1, $2, $3)"
	_, err = r.DB.Exec(query, login, email, password)
	if err != nil {
		return fmt.Errorf("error inserting new user: %w", err)
	}

	return nil
}

func (r *Repo) SignIn(login, password string) (int, int, error) {
	query := fmt.Sprintf("SELECT id,lastclick FROM users WHERE login = '%s' AND password = '%s'", login, password)

	rows, err := r.DB.Query(query)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var id int
	var lastclick int
	for rows.Next() {
		err := rows.Scan(&id, &lastclick)
		if err != nil {
			return 0, 0, err
		}
		return id, lastclick, nil
	}
	return 0, 0, fmt.Errorf("no student found with login=%s", login)

}
func (r *Repo) VerifyStudent(login, password string) error {
	_, _, err := r.SignIn(login, password)
	return err
}

//func (r *Repo) CardExists(name string, category string) error {
//	query := fmt.Sprintf("SELECT id,name,price,category FROM cards WHERE name='%s' AND category='%s'", name, category)
//	row := r.DB.QueryRow(query)
//	newID := 0
//	newName := ""
//	newPrice := 0
//	newCategory := ""
//
//	err := row.Scan(&newID, &newName, &newPrice, &newCategory)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//func (r *Repo) UpdateCard(name string, price int, category string) error {
//	query := fmt.Sprintf("UPDATE cards SET name='%s', price='%d', category='%s' WHERE name='%s' AND category='%s'", name, price, category, name, category)
//	_, err := r.DB.Exec(query)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
////
////func (r *Repo) RecreateDB(db string) error {
////	connstr := fmt.Sprintf("host=%s port=%s user=%s dbname=postgres password=%s sslmode=%s ", r.host, r.port, r.user, r.password, r.sslmode)
////	var err error
////	r.DB, err = sql.Open("postgres", connstr)
////	if err != nil {
////		return err
////	}
////	_, err = r.DB.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", db))
////	if err != nil {
////		return err
////	}
////	SetupDB()
////	return nil
////}
