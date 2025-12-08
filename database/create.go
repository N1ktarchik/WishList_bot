package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func ConnectToDB() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file. %w", err)
	}
	connStr := os.Getenv("CONNSTR")
	if connStr == "" {
		return nil, errors.New("connStr reading error")
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		err := fmt.Errorf("faild to open database. %w", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("faild connect to database. %w", err)
	}

	db.SetMaxOpenConns(30) // ~1000 сообщений/час
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)

	return db, nil

}

func CreateTables(db *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tables := []string{

		//chat id
		//username
		`CREATE TABLE IF NOT EXISTS users(
		id BIGINT PRIMARY KEY,
		username VARCHAR(255) UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		//numeric
		//chat id
		//link to users->chatID
		//name of wish
		//description
		//price
		//recerved (true/false)
		`CREATE TABLE IF NOT EXISTS wishes(
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(500) NOT NULL,
		description TEXT,
		link VARCHAR(1000),
		price DECIMAL(10,2),
		is_reserved BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,

		`CREATE INDEX IF NOT EXISTS idx_wishes_user_id ON wishes(user_id);`,
	}

	for _, query := range tables {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
	}

	for _, query := range indexes {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
	}

	return nil
}
