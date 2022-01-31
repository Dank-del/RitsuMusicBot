package database

type User struct {
	UserID         int64 `gorm:"primaryKey" json:"user_id"`
	LastFmUsername string
}

type BotUser struct {
	UserID      int64  `gorm:"primaryKey" json:"user_id"`
	UserName    string `json:"user_name"`
	ShowProfile bool   `json:"show_profile"`
}

type Chat struct {
	ChatID        int64  `gorm:"primaryKey"`
	StatusMessage string `gorm:"default:status"`
	DetectLinks   bool   `gorm:"default:true"`
}
