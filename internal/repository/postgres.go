package repository

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Добавлено для поддержки файлового источника
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	usersTable  = "users"
	pixelsTable = "pixels"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	// Подключение к базе "postgres" для проверки и создания основной БД
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=postgres password='%s' sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.SSLMode)
	initialDB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to initial database: %w", err)
	}
	defer initialDB.Close()

	// Проверяем наличие и создаем базу данных, если её нет
	if err := createDB(initialDB, cfg.DBName); err != nil {
		return nil, err
	}

	// Подключаемся к основной базе данных
	mainConnStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password='%s' sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	db, err := sqlx.Open("postgres", mainConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open main database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping main database: %w", err)
	}

	// Настраиваем миграции
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema", // Убедитесь, что путь правильный
		cfg.DBName, driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	logrus.Println("Migrations applied successfully!")
	return db, nil
}

func createDB(db *sqlx.DB, dbName string) error {
	// Проверяем наличие базы данных
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", dbName)
	logrus.Println(query)
	if err := db.Get(&exists, query); err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Создаем базу данных, если она не существует
	if !exists {
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		logrus.Printf("Database '%s' created successfully!", dbName)
	}

	return nil
}
