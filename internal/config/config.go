package config

import (
	"Market_backend/internal/auth"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

var Cfg Config
var AllowedOrigins []string

var (
	AppPort string

	RedisHost string
	RedisPort string

	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string

	S3Endpoint  string
	S3Host      string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    string
)

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

	AllowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	SMTPHost = os.Getenv("SMTP_HOST")
	SMTPPort = os.Getenv("SMTP_PORT")
	SMTPEmail = os.Getenv("SMTP_EMAIL")
	SMTPPassword = os.Getenv("SMTP_PASSWORD")

	S3Endpoint = os.Getenv("S3_ENDPOINT")
	S3Host = os.Getenv("S3_HOST")
	S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	S3SecretKey = os.Getenv("S3_SECRET_KEY")
	S3Bucket = os.Getenv("S3_BUCKET")
	S3UseSSL = os.Getenv("S3_USE_SSL")

	AppPort = os.Getenv("APP_PORT")
}
