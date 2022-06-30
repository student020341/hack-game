package db

import (
	"log"

	"server/pkg/accounts"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&accounts.Account{},
		&accounts.AuthSession{},
	)

	// create admin account if one does not exist
	var acc accounts.Account
	existingLogger := db.Logger
	db.Logger = logger.Discard
	tx := db.Take(&acc, "Level = 0")
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			acc, err := accounts.CreateAccount("admin", "admin")
			acc.Level = 0 // admin
			if err != nil {
				log.Fatalf("failed to create initial admin account: %+v", err)
			}
			tx = db.Create(&acc)
			if tx.Error != nil {
				log.Fatalf("failed to save initial admin account: %+v", tx.Error)
			}
		} else {
			log.Fatalf("error taking account on db init: %+v", tx.Error)
		}
	}

	db.Logger = existingLogger

	return db
}
