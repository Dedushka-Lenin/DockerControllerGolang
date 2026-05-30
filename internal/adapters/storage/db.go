package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
)

func GetDB(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)

	return sql.Open("postgres", psqlInfo)
}
