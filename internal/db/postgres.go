package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(dbURL string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	maxRetries := 15

	for i := 0; i < maxRetries; i++ {

		err = db.Ping()

		if err == nil {
			fmt.Println("Connected to PostgreSQL")
			return db, nil
		}

		fmt.Println("Database not ready, retrying...", err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("database not ready after retries: %w", err)
}
