package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ChatID    int64
	UserName  string
	CreatedAt time.Time
}

func (u *User) AddToDB(db *sql.DB) error {

	if u.ChatID <= 0 || u.UserName == "" {
		err := errors.New("data read error")
		return err
	}

	query :=
		`INSERT INTO users (id,username)
		 VALUES ($1,$2)
		 ON CONFLICT (id)
		 DO UPDATE SET username = EXCLUDED.username`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, u.ChatID, u.UserName)
	if err != nil {
		err := fmt.Errorf("error writing to the database. %w", err)
		return err
	}

	return nil
}

func GetUsernameByID(id int64, db *sql.DB) (string, error) {
	if id <= 0 {
		err := errors.New("wishID can not be empty")
		return "", err
	}

	query := `SELECT username FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var username string

	err := db.QueryRowContext(ctx, query, id).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("wish not found")
		}
		return "", fmt.Errorf("error reading wish from DB: %w", err)
	}

	if username == "" {
		return "", fmt.Errorf("user not faund")
	}

	return username, nil
}

func GetIdByUsername(username string, db *sql.DB) (int64, error) {
	if username == "" {
		err := errors.New("username can not be empty")
		return 0, err
	}

	username = strings.TrimPrefix(username, "@")

	query := `SELECT id FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64

	err := db.QueryRowContext(ctx, query, username).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("username not faund")
		}
		return 0, fmt.Errorf("error reading user from DB: %w", err)
	}

	if id == 0 {
		return 0, fmt.Errorf("user not faund")
	}

	return id, nil
}
