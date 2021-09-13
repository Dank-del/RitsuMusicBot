package config

type DaemonConfig struct {
	BotToken     string  `json:"bot_token"`
	LastFMKey    string  `json:"last_fm_key"`
	GeniusApiKey string  `json:"genius_api_key"`
	SudoUsers    []int64 `json:"sudo_users"`
}
