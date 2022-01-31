package config

import (
	"github.com/Dank-del/MusixScrape/musixScrape"
	"gorm.io/gorm"
)

type DaemonConfig struct {
	BotToken     string  `json:"bot_token"`
	LastFMKey    string  `json:"last_fm_key"`
	SudoUsers    []int64 `json:"sudo_users"`
}

type DaemonLocal struct {
	Config            *DaemonConfig
	SqlSession        *gorm.DB
	MusixMatchSession *musixScrape.Client
}
