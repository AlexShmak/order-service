package models

import "time"

type User struct {
	ID           int64
	Name         string
	PasswordHash []byte
	Email        string
	CreatedAt    time.Time
}
