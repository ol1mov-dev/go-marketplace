package databases

import (
	"database/sql"
	"fmt"
	"log"
	"order-service/configs"
	"time"

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

func ConnectPostgreSQL() error {
	log.Println("Loading PostgreSQL config...")
	dbConfig, err := configs.LoadPgConfig()

	log.Println("Making connection string...")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName, dbConfig.SSLMode)

	log.Println("Opening connection...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {

		return fmt.Errorf("failed to open databases: %v", err)
	}

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to databases: %v", err)
	}

	log.Println("Successfully connected to databases!")

	PostgresDB = db

	return nil
}
