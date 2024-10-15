package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Port     string
}
type DatabaseConfig struct {
	User     string
	Password string
	DBName   string
}

func Load() Config {
	godotenv.Load(".env")
	return Config{
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		Port: os.Getenv("DB_PORT"),
	}
}
