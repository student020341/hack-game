package models

import "time"

type Account struct {
	ID       string `gorm:"unique;not null"`
	Username string `gorm:"unique;not null"`
	Salt     []byte `json:"-"`
	Password string `json:"-" gorm:"not null"`
	// 0 = admin, 1 = mod, 2 = user?
	Level      int
	Characters []Character
}

type AuthSession struct {
	ID           string `gorm:"not null"`
	Token        string `gorm:"not null"`
	CreatedAt    time.Time
	AccountID    string `gorm:"not null"`
	LastAccessed time.Time
}

type Character struct {
	ID        string `gorm:"not null"`
	AccountID string
	Name      string
	// TODO level/exp?
	Inventory []Item
}

type Item struct {
	ID          string `gorm:"not null"`
	Something   string // TODO
	CharacterID string
}
