package database

import (
	"fmt"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gorm.io/gorm/logger"

	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDatabase(botId int64) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%d.db", botId)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logging.SUGARED.Error("failed to connect database")
	}
	config.Local.SqlSession = db
	logging.SUGARED.Info("Database connected")

	// Create tables if they don't exist
	err = config.Local.SqlSession.AutoMigrate(&User{}, &BotUser{}, &Chat{})
	if err != nil {
		logging.SUGARED.Error(err)
	}

	logging.SUGARED.Info("Auto-migrated database schema")

}
