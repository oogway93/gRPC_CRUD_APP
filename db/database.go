package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func DatabaseConnection(cfg Config) *gorm.DB {
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
	db.Exec("CREATE SCHEMA IF NOT EXISTS public;")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
