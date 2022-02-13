package config

import (
	"github.com/Dank-del/MusixScrape/musixScrape"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"gitlab.com/toby3d/telegraph"
	"gorm.io/gorm"
)

type DaemonConfig struct {
	BotToken            string  `json:"bot_token"`
	LastFMKey           string  `json:"last_fm_key"`
	SudoUsers           []int64 `json:"sudo_users"`
	SpotifyClientID     string  `json:"spotify_client_id"`
	SpotifyClientSecret string  `json:"spotify_client_secret"`
	SpotifyRedirectUri  string  `json:"spotify_redirect_uri"`
	ServerAddr          string  `json:"server_addr"`
}

type DaemonLocal struct {
	Config            *DaemonConfig
	SqlSession        *gorm.DB
	MusixMatchSession *musixScrape.Client
	TelegraphClient   *telegraph.Account
	Bot               *gotgbot.Bot
}
