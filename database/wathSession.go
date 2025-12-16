package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type WishSession struct {
	ChatID   int64
	TargetID int64
	WishID   int64
	live     time.Time
}

// wishSession
func (ws *WishSession) IsAlive() bool {
	return time.Now().Before(ws.live)
}

func (ws *WishSession) UpdateLiveTime(minutes int) {
	ws.live = time.Now().Add(time.Duration(minutes) * time.Minute)
}

func CreateNewWishSession(chatid, target int64, db *sql.DB) (*WishSession, error) {
	newSession := WishSession{
		ChatID:   chatid,
		TargetID: target,
		live:     time.Now().Add(30 * time.Minute),
	}
	query := `SELECT MIN(id) FROM wishes WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := db.QueryRowContext(ctx, query, target).Scan(&newSession.WishID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения минимального ID: %w", err)
	}

	err = newSession.Save(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения новой сессии: %w", err)
	}
	return &newSession, nil

}

func (ws *WishSession) Reset() {
	ws.WishID = 0
	ws.TargetID = 0
	ws.UpdateLiveTime(30)
}

func (ws *WishSession) Update(db *sql.DB) error {
	if ws.WishID == 0 || ws.ChatID == 0 || ws.TargetID == 0 {
		err := errors.New("data read error")
		return err
	}

	query := `
        UPDATE watch_session SET current_wish = $1 , live=$2 WHERE chat_id = $3 `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query,
		ws.WishID,
		ws.live,
		ws.ChatID,
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

func (ws *WishSession) Save(db *sql.DB) error {
	query := `INSERT INTO watch_session(id,current_wish,target_id,live)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT(id) DO UPDATE SET 
				current_wish = EXCLUDED.current_wish ,
				target_id = EXCLUDED.target_id ,
				live = EXCLUDED.live`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := db.ExecContext(ctx, query, ws.ChatID, ws.WishID, ws.TargetID, ws.live)

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

func (ws *WishSession) Delete(db *sql.DB) error {

	if ws.ChatID <= 0 {
		err := errors.New("ChatID can not be empty")
		return err
	}

	query := `DELETE FROM watch_session WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.ExecContext(ctx, query, ws.ChatID)

	if err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			err := fmt.Errorf("delete by timeout. %w", err) //не кончился ли таймаут
			return err
		}

		err := fmt.Errorf("error deleting wish. %w", err)
		return err
	}

	ws.Reset()
	return nil
}

func GetWishSessonByID(chatID int64, db *sql.DB) (*WishSession, error) {
	if chatID <= 0 {
		err := fmt.Errorf("chatID can not be zero")
		return nil, err
	}
	var session WishSession

	query := `SELECT * FROM watch_session WHERE id =$1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := db.QueryRowContext(ctx, query, chatID).Scan(&session.ChatID, &session.TargetID, &session.WishID, &session.live)
	if err != nil {
		return nil, fmt.Errorf("get session from DB error. %w", err)
	}

	if !session.IsAlive() {
		session.Delete(db)
		return nil, nil
	}

	return &session, nil

}

func CleanExpiredSessions(db *sql.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	query := `DELETE FROM watch_session WHERE live < NOW()`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("cleanup sessions error: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.Printf("Cleaned up %d expired sessions", rows)
	}

	return nil
}
