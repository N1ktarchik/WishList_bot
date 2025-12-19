package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type WishNavigation struct {
	NextID, PrevID *int64
}

func GetWishNavigation(currentWishID, targetID int64, db *sql.DB) (*WishNavigation, error) {
	if currentWishID <= 0 || targetID <= 0 {
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

	err := db.QueryRowContext(ctx, query, targetID, currentWishID).Scan(&prev, &next)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
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
