package models

import "time"

type Account struct {
	ID       string
	Username string `gorm:"unique;not null"`
	Salt     []byte
	Password string `gorm:"not null"`
	// 0 = admin, 1 = mod, 2 = user?
	Level int
}

type AuthSession struct {
	ID           string
	Token        string
	CreatedAt    time.Time
	AccountID    string
	LastAccessed time.Time
}
