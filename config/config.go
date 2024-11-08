package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func NewCfg() (*DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return &DB{}, fmt.Errorf("error loading .env file. Please create it in the root of project: %v", err)
	}
	db := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	return db, nil
}
