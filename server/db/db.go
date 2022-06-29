package db

import (
	"fmt"

	"server/pkg/accounts"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDB(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&accounts.Account{},
	)

	return db
}

func Init() {
	DB = NewDB("test.db")

	fmt.Println("connected to db")
}
