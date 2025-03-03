package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf(" DATABASE_URL is not set in environment variables")
	}

	var err error
	for i := 0; i < 10; i++ { // Retry 10 times
		DB, err = sql.Open("postgres", dbURL)
		if err == nil && DB.Ping() == nil {
			fmt.Println(" Connected to PostgreSQL successfully")
			return nil
		}

		fmt.Println("Database not ready yet. Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after multiple retries: %v", err)
}
