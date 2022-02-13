package database

import (
	"errors"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/zmb3/spotify/v2"
	"gitlab.com/Dank-del/lastfm-tgbot/core/auth"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"strings"

	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
)

func UpdateChat(ChatId int64, statusMessage string) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := chatsMap[ChatId]
	if data != nil && data.ChatID == ChatId && data.StatusMessage == statusMessage {
		return
	}
	tx := config.Local.SqlSession.Begin()
	chat := &Chat{
		ChatID:        ChatId,
		StatusMessage: statusMessage,
	}
	chatsMap[ChatId] = chat
	tx.Save(chat)
	tx.Commit()
}

func GetChat(chatID int64) (c *Chat, err error) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := chatsMap[chatID]
	if data != nil && data.ChatID == chatID {
		return data, nil
	}
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := Chat{}
	config.Local.SqlSession.Where("chat_id = ?", chatID).Take(&p)
	chatsMap[chatID] = &p
	return &p, nil
}

func UpdateLastFMUserInDB(UserID int64, LastFmUsername string) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := usersMap[UserID]
	if data != nil && data.LastFmUsername == LastFmUsername {
		return
	}
	tx := config.Local.SqlSession.Begin()
	user := &User{
		UserID:         UserID,
		LastFmUsername: strings.ToLower(LastFmUsername),
	}
	usersMap[UserID] = user
	tx.Save(user)
	tx.Commit()
}

func UpdateBotUser(UserID int64, UserName string, ShowProfile bool) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := botUserMapById[UserID]
	if data != nil && data.UserID == UserID && data.UserName == UserName && data.ShowProfile == ShowProfile {
		return
	}
	tx := config.Local.SqlSession.Begin()
	user := &BotUser{
		UserID:      UserID,
		UserName:    UserName,
		ShowProfile: ShowProfile,
	}
	botUserMapById[UserID] = user
	tx.Save(user)
	tx.Commit()
}

func GetLastFMUserFromDB(UserID int64) (u *User, err error) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := usersMap[UserID]
	if data != nil && data.UserID == UserID {
		return data, nil
	}
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := User{}
	config.Local.SqlSession.Where("user_id = ?", UserID).Take(&p)
	usersMap[UserID] = &p
	return &p, nil
}

func GetBotUserByID(UserID int64) (u *BotUser, err error) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := botUserMapById[UserID]
	if data != nil && data.UserID == UserID {
		return data, nil
	}
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	config.Local.SqlSession.Where("user_id = ?", UserID).Take(&p)
	botUserMapById[UserID] = &p
	return &p, nil
}

func GetBotUserByUsername(UserName string) (u *BotUser, err error) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := botUserMapByUsername[UserName]
	if data != nil && data.UserName == UserName {
		return data, nil
	}
	if config.Local.SqlSession == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := BotUser{}
	config.Local.SqlSession.Where("user_name = ?", UserName).Take(&p)
	botUserMapByUsername[UserName] = &p
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

func UpdateSpotifyUser(userId int64, token string) {
	var usr *auth.SpotifyUser
	if usr = spotifyUserMap[userId]; usr != nil && usr.RefreshToken == token {
		return
	}
	ctx := context.Background()
	userAuth := spotify.New(auth.SpotifyAuthenticator.Client(ctx, &oauth2.Token{RefreshToken: token}))
	_, err := userAuth.CurrentUser(ctx)
	if err != nil {
		return
	}
	usr = &auth.SpotifyUser{
		UserId:       userId,
		RefreshToken: token,
	}
	spotifyUserMap[userId] = usr
	tx := config.Local.SqlSession.Begin()
	tx.Save(usr)
	tx.Commit()
}

func GetSpotifyUser(userId int64) (*spotify.Client, error) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	ctx := context.Background()
	var usr *auth.SpotifyUser
	usr = spotifyUserMap[userId]
	if usr != nil && usr.UserId == userId {
		userAuth := spotify.New(auth.SpotifyAuthenticator.Client(ctx, &oauth2.Token{RefreshToken: usr.RefreshToken}))
		_, err := userAuth.CurrentUser(ctx)
		if err != nil {
			return nil, err
		}
		return userAuth, nil
	}
	config.Local.SqlSession.Where(&auth.SpotifyUser{UserId: userId}).Take(&usr)
	if usr == nil || usr.RefreshToken == strongStringGo.EMPTY {
		return nil, errors.New("user not found")
	}
	userAuth := spotify.New(auth.SpotifyAuthenticator.Client(ctx, &oauth2.Token{RefreshToken: usr.RefreshToken}))
	_, err := userAuth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	spotifyUserMap[userId] = usr
	return userAuth, nil
}
