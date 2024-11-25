package main

import (
	"fmt"
	dewu "github.com/Xapsiel/PBCFU"
	"github.com/Xapsiel/PBCFU/internal/handler"
	"github.com/Xapsiel/PBCFU/internal/repository"
	"github.com/Xapsiel/PBCFU/internal/service"
	"github.com/Xapsiel/PBCFU/internal/service/log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

func main() {
	output, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY, 0666)
	log.Logger = log.NewLogService(output)
	if err != nil {
		fmt.Errorf("Ошибка логирования")
		return
	}
	if err := initConfig(); err != nil {
		log.Logger.Warn(-1, err.Error())
	}
	if err := godotenv.Load(); err != nil {

		log.Logger.Warn(-1, "Error loading .env file")
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Logger.Warn(-1, err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos, output)
	handlers := handler.NewHandler(services)

	srv := new(dewu.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Logger.Warn(-1, err.Error())
	}

}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
