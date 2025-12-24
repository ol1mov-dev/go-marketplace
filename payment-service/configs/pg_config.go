package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func LoadPgConfig() (*PgConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return &PgConfig{}, err
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return &PgConfig{host, port, user, password, database, sslmode}, nil
}
