package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set in environment variables")
	}

	var err error
	for i := 0; i < 10; i++ { // Retry 10 times
		DB, err = pgxpool.New(context.Background(), databaseURL)
		if err == nil {
			err = DB.Ping(context.Background())
			if err == nil {
				fmt.Println("Connected to PostgreSQL successfully")
				return nil
			}
		}

		fmt.Println("Database not ready yet. Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after multiple retries: %v", err)
}
