package database

import (
	"database/sql"
	"fmt"

	"github.com/todanni/auth/config"
)

func Open(cfg config.Config) (*sql.DB, error) {
	// Make connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	//db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})

	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}
