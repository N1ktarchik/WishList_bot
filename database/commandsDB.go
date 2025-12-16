package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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

func GetWishNavigation(currentWishID, chatID int64, db *sql.DB) (*WishNavigation, error) {
	if currentWishID <= 0 || chatID <= 0 {
		err := errors.New("ivalid ID")
		return nil, err
	}

	query := `
	SELECT MAX(CASE WHEN id<$2 THEN id END),
		MIN(CASE WHEN id>$2 THEN id END)
	FROM wishes WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var next, prev sql.NullInt64

	err := db.QueryRowContext(ctx, query, chatID, currentWishID).Scan(&prev, &next)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wish or user not found")
		}
		return nil, fmt.Errorf("error reading data from DB: %w", err)
	}

	result := &WishNavigation{}

	if next.Valid {
		result.NextID = &next.Int64
	}

	if prev.Valid {
		result.PrevID = &prev.Int64
	}

	return result, nil
}

func (w *Wish) UpdateWish(db *sql.DB) error { //

	if w.ChatIdLink == 0 || w.WishName == "" || w.Price < 0 || w.ID <= 0 {
		err := errors.New("data read error")
		return err
	}

	query := `
        UPDATE wishes SET name = $1, description = $2, link = $3, price = $4, is_reserved = $5 WHERE id = $6 `

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
		w.IsReserved,
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

func CleanOverdueStatuses(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `DELETE FROM user_status WHERE live<NOW()`

	_, err := db.ExecContext(ctx, query)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error delete dead states. %w", err)
		return err
	}

	return nil
}

func (s *UserStatus) Save(db *sql.DB) error {

	query := `INSERT INTO user_status(id,step,name,description,link,price,live)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT(id) DO UPDATE SET 
				step = EXCLUDED.step ,
				name = EXCLUDED.name ,
				description = EXCLUDED.description ,
				link = EXCLUDED.link , 
				price = EXCLUDED.price ,
				live = EXCLUDED.live`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var dbPrice interface{}
	if s.Price == 0 {
		dbPrice = nil
	} else {
		dbPrice = s.Price
	}

	_, err := db.ExecContext(ctx, query, s.ChatID, s.Step, s.WishName, s.Description, s.Url, dbPrice, s.live)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("cancellation by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error updating user status. %w", err)
		return err
	}

	return nil
}

func (s *UserStatus) Delete(db *sql.DB) error {

	if s.ChatID <= 0 {
		err := errors.New("ChatID can not be empty")
		return err
	}

	query := `DELETE FROM user_status WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.ExecContext(ctx, query, s.ChatID)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error deleting wish. %w", err)
		return err
	}

	s.Reset()
	return nil
}

func GetUserStatusByID(db *sql.DB, id int64) (*UserStatus, error) {
	if id <= 0 {
		err := errors.New("ChatID can not be empty")
		return nil, err
	}

	query := `SELECT * FROM user_status WHERE id = $1 AND live>NOW()`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	row := db.QueryRowContext(ctx, query, id)

	status := &UserStatus{}
	var price sql.NullFloat64

	err := row.Scan(
		&status.ChatID,
		&status.Step,
		&status.WishName,
		&status.Description,
		&status.Url,
		&price,
		&status.live,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		err := fmt.Errorf("Get user status error. %w", err)
		return nil, err
	}

	if price.Valid {
		status.Price = price.Float64
	} else {
		status.Price = 0
	}

	return status, nil

}

func (s *UserStatus) IsAlive() bool {
	return time.Now().Before(s.live)
}

func (s *UserStatus) UpdateLiveTime(minutes int) {
	s.live = time.Now().Add(time.Duration(minutes) * time.Minute)
}

func (s *UserStatus) Reset() {
	s.Step = 0
	s.WishName = ""
	s.Description = ""
	s.Url = ""
	s.Price = 0
	s.UpdateLiveTime(10)
}

func CreateNewUserStatus(id int64) *UserStatus {
	return &UserStatus{
		ChatID: id,
		Step:   1,
		Price:  0,
		live:   time.Now().Add(10 * time.Minute),
	}
}
