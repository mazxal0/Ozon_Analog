package config

import (
	"Market_backend/internal/auth"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

var Cfg Config

func Init() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Println("DATABASE_URL not set, using default")
		dbUrl = "postgres://eduvix:eduvix_pass@localhost:5432/eduvix_db?sslmode=disable"
	}

	secret := os.Getenv("JWT_SECRET")
	auth.InitJwt(secret)

	Cfg.DBUrl = dbUrl
}
