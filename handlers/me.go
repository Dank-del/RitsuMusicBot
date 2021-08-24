package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"html"
)

func meHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	getusername, err := database.GetUser(msg.From.Id)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", html.EscapeString(err.Error())), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	if getusername == nil {
		_, err := msg.Reply(b, "<i>Error: lastfm username not set, use /setusername</i>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	lastFMuser, _ := last_fm.GetLastFMUser(getusername.LastFmUsername)
	m := fmt.Sprintf("<b>%s</b>\n\n", lastFMuser.User.Name)
	m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)
	m += fmt.Sprintf("<b>Playlist Count</b>: %s\n", lastFMuser.User.Playlists)
	m += fmt.Sprintf("<b>Gender</b>: %s\n", lastFMuser.User.Gender)
	m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)
	m += fmt.Sprintf("<b>Age</b>: %s\n", lastFMuser.User.Age)

	_, err = b.SendPhoto(msg.Chat.Id, lastFMuser.User.Image[3].Text, &gotgbot.SendPhotoOpts{ParseMode: "html", Caption: m,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "View on Last.FM", Url: lastFMuser.User.URL},
		}}}})
	return err
}
