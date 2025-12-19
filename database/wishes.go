package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Wish struct {
	ID                         int64
	ChatIdLink                 int64
	WishName, Description, Url string
	Price                      float64
	IsReserved                 bool
	CreatedAt                  time.Time
}

func (w *Wish) AddToDB(db *sql.DB) error { //

	if w.ChatIdLink <= 0 || w.WishName == "" || w.Price < 0 {
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

	var dbPrice interface{}
	if w.Price == 0 {
		dbPrice = nil
	} else {
		dbPrice = w.Price
	}

	err := db.QueryRowContext(ctx, query,
		w.ChatIdLink,
		w.WishName,
		w.Description,
		w.Url,
		dbPrice).Scan(
		&w.ID,
		&w.IsReserved,
		&w.CreatedAt,
	)

	return err

}

func (w *Wish) DeleteFromDB(db *sql.DB) error {

	if w.ID <= 0 {
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

	rows, err := result.RowsAffected() //RowsAfected сколько строк измененно
	if err != nil {
		return fmt.Errorf("error checking deletion: %w", err)
	}

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
	var price sql.NullFloat64
	err := db.QueryRowContext(ctx, query, id).Scan(
		&wish.ID,
		&wish.ChatIdLink,
		&wish.WishName,
		&wish.Description,
		&wish.Url,
		&price,
		&wish.IsReserved,
		&wish.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wish not found")
		}
		return nil, fmt.Errorf("error reading wish from DB: %w", err)
	}

	if price.Valid {
		wish.Price = price.Float64
	} else {
		wish.Price = 0
	}

	return &wish, nil
}

func (w *Wish) UpdateWish(db *sql.DB) error { //

	if w.ChatIdLink == 0 || w.WishName == "" || w.Price < 0 || w.ID <= 0 {
		err := errors.New("data read error")
		return err
	}

	query := `
        UPDATE wishes SET name = $1, description = $2, link = $3, price = $4 WHERE id = $5 `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbPrice interface{}
	if w.Price == 0 {
		dbPrice = nil
	} else {
		dbPrice = w.Price
	}

	result, err := db.ExecContext(ctx, query,
		w.WishName,
		w.Description,
		w.Url,
		dbPrice,
		w.ID,
	)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error updating wish. %w", err)
		return err
	}

	rows, err := result.RowsAffected() //RowsAfected сколько строк измененно
	if err != nil {
		return fmt.Errorf("error checking updating: %w", err)
	}

	if rows == 0 {
		err := errors.New("wish not found")
		return err
	}

	return nil

}

func ReserveWish(wishID int64, db *sql.DB) error {
	if wishID <= 0 {
		err := errors.New("id can not be zero")
		return err
	}

	query := `UPDATE wishes SET is_reserved = $1 WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query, true, wishID)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error recerving wish. %w", err)
		return err
	}

	rows, err := result.RowsAffected() //RowsAfected сколько строк измененно
	if err != nil {
		return fmt.Errorf("error checking updating: %w", err)
	}

	if rows == 0 {
		err := errors.New("wish not found")
		return err
	}

	return nil
}
