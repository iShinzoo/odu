package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(dbURL string) (*sql.DB, error) {

	var db *sql.DB
	var err error

	maxRetries := 10

	for i := 0; i < maxRetries; i++ {

		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return nil, err
		}

		// connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		err = db.Ping()

		if err == nil {
			fmt.Println("Connected to PostgreSQL")
			return db, nil
		}

		fmt.Println("Database not ready, retrying...")
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts", maxRetries)
}
