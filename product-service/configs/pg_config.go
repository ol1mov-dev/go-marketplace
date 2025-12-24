package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type PgConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
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
	maxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	connMaxLifetime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	connMaxIdleTime, err := time.ParseDuration(os.Getenv("DB_CONN_MAX_IDLE_TIME"))

	return &PgConfig{
		host,
		port,
		user,
		password,
		database,
		sslmode,
		maxOpenConns,
		maxIdleConns,
		connMaxLifetime,
		connMaxIdleTime,
	}, nil
}
