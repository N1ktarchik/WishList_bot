package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserStatus struct {
	ChatID                     int64
	Step                       int
	WishName, Description, Url string
	Price                      float64
	live                       time.Time
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
