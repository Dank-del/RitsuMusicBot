package database

import (
	"fmt"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var SESSION *gorm.DB

func StartDatabase(botId int64) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%d.db", botId)), &gorm.Config{})
	if err != nil {
		logging.Error("failed to connect database")
	}

	SESSION = db
	logging.Info("Database connected")

	// Create tables if they don't exist
	err = SESSION.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
	}

	logging.Info("Auto-migrated database schema")

}
