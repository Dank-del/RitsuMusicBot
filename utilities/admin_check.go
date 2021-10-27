package utilities

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func contains(s []gotgbot.ChatMember, e int64) bool {
	for _, a := range s {
		if a.GetUser().Id == e {
			return true
		}
	}
	return false
}

func IsUserAdmin(bot *gotgbot.Bot, chat *gotgbot.Chat, userId int64) bool {
	if chat.Type == "private" {
		return true
	}

	var adminList []gotgbot.ChatMember

	adminList, err := chat.GetAdministrators(bot)
	if err != nil { // not found
		return false
	}
	return contains(adminList, userId)
}
