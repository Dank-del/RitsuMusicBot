package database

import (
	"errors"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"strings"
)

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
}

func UpdateChat(ChatId int64, statusMessage string) {
	tx := config.Local.SqlSession.Begin()
	chat := &Chat{ChatID: ChatId, StatusMessage: statusMessage}
	tx.Save(chat)
	tx.Commit()
}

func GetChat(chatID int64) (c *Chat, err error) {
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := Chat{}
	config.Local.SqlSession.Where("chat_id = ?", chatID).Take(&p)
	return &p, nil
}

func (c *Chat) GetStatusMessage() string {
	if c == nil || c.StatusMessage == "" {
		return "status"
	}
	return c.StatusMessage
}

func UpdateLastFMUserInDB(UserID int64, LastFmUsername string) {
	tx := config.Local.SqlSession.Begin()
	user := &User{UserID: UserID, LastFmUsername: strings.ToLower(LastFmUsername)}
	tx.Save(user)
	tx.Commit()
}

func UpdateBotUser(UserID int64, UserName string, ShowProfile bool) {
	tx := config.Local.SqlSession.Begin()
	user := BotUser{UserID: UserID, UserName: UserName, ShowProfile: ShowProfile}
	tx.Save(user)
	tx.Commit()
}

func GetLastFMUserFromDB(UserID int64) (u *User, err error) {
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := User{}
	config.Local.SqlSession.Where("user_id = ?", UserID).Take(&p)
	return &p, nil
}

func GetBotUserByID(UserID int64) (u *BotUser, err error) {
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	config.Local.SqlSession.Where("user_id = ?", UserID).Take(&p)
	return &p, nil
}

func GetBotUserByUsername(UserName string) (u *BotUser, err error) {

	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	config.Local.SqlSession.Where("user_name = ?", UserName).Take(&p)
	return &p, nil

}

func GetBotUserCount() (c int64) {
	config.Local.SqlSession.Model(&BotUser{}).Count(&c)
	return c
}

func GetLastmUserCount() (c int64) {
	config.Local.SqlSession.Model(&User{}).Count(&c)
	return c
}
