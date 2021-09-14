package database

import (
	"fmt"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var SESSION *gorm.DB

func StartDatabase(botId int64) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%d.db", botId)), &gorm.Config{})
	if err != nil {
		logging.SUGARED.Error("failed to connect database")
	}

	SESSION = db
	logging.SUGARED.Info("Database connected")

	// Create tables if they don't exist
	err = SESSION.AutoMigrate(&User{}, &BotUser{})
	if err != nil {
		logging.SUGARED.Error(err)
	}

	logging.SUGARED.Info("Auto-migrated database schema")

}
