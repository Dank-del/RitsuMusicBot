package database

import "gitlab.com/Dank-del/lastfm-tgbot/core/config"

func (u *User) IsValid() bool {
	return u != nil && u.UserID != 0 && u.LastFmUsername != ""
}

func (c *Chat) GetStatusMessage() string {
	if c == nil || c.StatusMessage == "" {
		return "status"
	}
	return c.StatusMessage
}

func (c *Chat) SetLinkDetection(enabled bool) {
	databaseMutex.RLock()
	defer databaseMutex.RUnlock()
	data := chatsMap[c.ChatID]
	if data != nil && data.DetectLinks == enabled {
		return
	}
	tx := config.Local.SqlSession.Begin()
	c.DetectLinks = enabled
	chatsMap[c.ChatID] = c
	tx.Save(c)
	tx.Commit()
}
