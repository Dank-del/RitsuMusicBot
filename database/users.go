package database

import (
	"errors"
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
	tx := SESSION.Begin()
	chat := &Chat{ChatID: ChatId, StatusMessage: statusMessage}
	tx.Save(chat)
	tx.Commit()
}

func GetChat(chatID int64) (c *Chat, err error) {
	if SESSION == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := Chat{}
	SESSION.Where("chat_id = ?", chatID).Take(&p)
	return &p, nil
}

func (c *Chat) GetStatusMessage() string {
	if c == nil || c.StatusMessage == "" {
		return "status"
	}
	return c.StatusMessage
}

func UpdateLastFMUserInDB(UserID int64, LastFmUsername string) {
	tx := SESSION.Begin()
	user := &User{UserID: UserID, LastFmUsername: strings.ToLower(LastFmUsername)}
	tx.Save(user)
	tx.Commit()
}

func UpdateBotUser(UserID int64, UserName string, ShowProfile bool) {
	tx := SESSION.Begin()
	user := BotUser{UserID: UserID, UserName: UserName, ShowProfile: ShowProfile}
	tx.Save(user)
	tx.Commit()
}

func GetLastFMUserFromDB(UserID int64) (u *User, err error) {
	if SESSION == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := User{}
	SESSION.Where("user_id = ?", UserID).Take(&p)
	return &p, nil
}

func GetBotUserByID(UserID int64) (u *BotUser, err error) {
	if SESSION == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	SESSION.Where("user_id = ?", UserID).Take(&p)
	return &p, nil
}

func GetBotUserByUsername(UserName string) (u *BotUser, err error) {

	if SESSION == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	SESSION.Where("user_name = ?", UserName).Take(&p)
	return &p, nil

}

func GetBotUserCount() (c int64) {
	SESSION.Model(&BotUser{}).Count(&c)
	return c
}

func GetLastmUserCount() (c int64) {
	SESSION.Model(&User{}).Count(&c)
	return c
}
