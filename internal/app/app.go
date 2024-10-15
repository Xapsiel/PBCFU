package app

import (
	"dewu/internal/config"
	"dewu/internal/database/postgresql"
	"dewu/internal/transport/rest"
)

func Run() error {
	cfg := config.Load()
	db, err := postgresql.New(cfg.Database)
	if err != nil {
		return err
	}
	return rest.StartServer(cfg, db.DB)
}
