package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"html"
	"time"
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
	if getusername.LastFmUsername == "A" {
		_, err := msg.Reply(b, "<i>Error: lastfm username not set, use /setusername</i>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	lastFMuser, _ := last_fm.GetLastFMUser(getusername.LastFmUsername)
	createdAt := time.Unix(int64(lastFMuser.User.Registered.Text), 0).Format(time.RFC850)

	m := fmt.Sprintf("<b>%s</b>\n\n", lastFMuser.User.Name)
	m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)
	m += fmt.Sprintf("<b>Playlist Count</b>: %s\n", lastFMuser.User.Playlists)
	m += fmt.Sprintf("<b>Gender</b>: %s\n", lastFMuser.User.Gender)
	// m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)
	m += fmt.Sprintf("<b>Age</b>: %s\n\n", lastFMuser.User.Age)
	m += fmt.Sprintf("<b>Created at</b>\n <code>%s</code>", createdAt)

	pic := lastFMuser.User.Image
	if pic == nil {
		_, err = msg.Reply(b, m, &gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}
	_, err = b.SendPhoto(msg.Chat.Id, pic[3].Text, &gotgbot.SendPhotoOpts{ParseMode: "html", Caption: m,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "View on Last.FM", Url: lastFMuser.User.URL},
		}}}})
	return err
}
