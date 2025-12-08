package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type commands interface {
	AddToDB(db *sql.DB) error
}

type User struct {
	ChatID    int64
	UserName  string
	CreatedAt time.Time
}

type Wish struct {
	ID                         int64
	ChatIdLink                 int64
	WishName, Description, Url string
	Price                      float64
	IsReserved                 bool
	CreatedAt                  time.Time
}

func (u User) AddToDB(db *sql.DB) error {

	if u.ChatID == 0 || u.UserName == "" {
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

func (w *Wish) AddToDB(db *sql.DB) error {

	if w.ChatIdLink == 0 || w.WishName == "" || w.Price < 0 {
		err := errors.New("data read error")
		return err
	}

	query := `
        INSERT INTO wishes (user_id, name, description, link, price)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, is_reserved, created_at
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query,
		w.ChatIdLink,
		w.WishName,
		w.Description,
		w.Url,
		w.Price).Scan(
		&w.ID,
		&w.IsReserved,
		&w.CreatedAt,
	)

	return err

}

func (w Wish) DeleteFromDB(db *sql.DB) error {

	if w.ID == 0 {
		err := errors.New("wishID can not be empty")
		return err
	}

	query := `DELETE FROM wishes WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query, w.ID)
	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error deleting wish. %w", err)
		return err
	}

	rows, _ := result.RowsAffected() //RowsAfected сколько строк измененно
	if rows == 0 {
		err := errors.New("wish not found")
		return err
	}

	return nil
}

func GetWishByID(id int64, db *sql.DB) (*Wish, error) {

	if id <= 0 {
		err := errors.New("wishID can not be empty")
		return nil, err
	}

	query := `SELECT * FROM wishes WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wish Wish
	err := db.QueryRowContext(ctx, query, id).Scan(
		&wish.ID,
		&wish.ChatIdLink,
		&wish.WishName,
		&wish.Description,
		&wish.Url,
		&wish.Price,
		&wish.IsReserved,
		&wish.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wish not found")
		}
		return nil, fmt.Errorf("error reading wish from DB: %w", err)
	}

	return &wish, nil
}
