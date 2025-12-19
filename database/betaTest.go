package database

//for beta-test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func SaveNewTester(chatID int64, username string, db *sql.DB) error {
	if chatID <= 0 {
		return errors.New("chatId can not be zero")
	}

	query := `INSERT INTO testers(id,username) VALUES($1,$2) ON CONFLICT(id) DO UPDATE SET username = EXCLUDED.username `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.ExecContext(ctx, query, chatID, username)
	if err != nil {
		return fmt.Errorf("error to save new tester: %w", err)
	}

	return nil
}

func ChekTesterRights(chatID int64, db *sql.DB) bool {

	var result int64

	query := `SELECT id FROM testers WHERE id =$1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := db.QueryRowContext(ctx, query, chatID).Scan(&result)

	if err != nil {
		return false
	}

	if result > 0 {
		return true
	}

	return false
}
