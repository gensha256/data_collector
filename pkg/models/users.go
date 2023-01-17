package models

import "time"

type User struct {
	ID        string
	Email     string
	Telegram  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
