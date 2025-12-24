package databases

import (
	"database/sql"
	"fmt"
	"log"
	"payment-service/configs"

	_ "github.com/lib/pq"
)

func ConnectPostgreSQL() (*sql.DB, error) {
	log.Println("Loading PostgreSQL config...")
	dbConfig, err := configs.LoadPgConfig()

	log.Println("Making connection string...")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName, dbConfig.SSLMode)

	log.Println("Opening connection...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {

		return nil, fmt.Errorf("failed to open databases: %v", err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to databases: %v", err)
	}

	log.Println("Successfully connected to databases!")
	return db, nil
}
