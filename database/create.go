package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	db, err := sql.Open("pgx", connStr)

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

		//chat id
		//step status
		//wish data
		//...
		//live time
		`CREATE TABLE IF NOT EXISTS user_status(
		id BIGINT PRIMARY KEY,
		step INTEGER DEFAULT 0,
		name VARCHAR(500) NOT NULL,
		description TEXT,
		link VARCHAR(1000),
		price DECIMAL(10,2),
		new BOOLEAN DEFAULT TRUE,
		live TIMESTAMP NOT NULL
		);`,

		//chat id
		//target chat id
		//wish id (numeric)
		//live time
		`CREATE TABLE IF NOT EXISTS watch_session(
		id BIGINT PRIMARY KEY,
		target_id BIGINT,
		current_wish BIGINT,
		live TIMESTAMP NOT NULL
		);`,
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,

		`CREATE INDEX IF NOT EXISTS idx_wishes_user_id ON wishes(user_id);`,

		`CREATE INDEX IF NOT EXISTS idx_watch_session_live ON watch_session(live);`,

		`CREATE INDEX IF NOT EXISTS idx_user_status_live ON user_status(live);`,
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
