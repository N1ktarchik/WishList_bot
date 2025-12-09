package database

import "time"

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

type WishNavigation struct {
	NextID, PrevID *int64
}
