package handlers

import (
	"fmt"
	"html"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
)

func meHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	getusername, err := database.GetLastFMUserFromDB(msg.From.Id)
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
	lastFMuser, _ := lastfm.GetLastFMUser(getusername.LastFmUsername)
	createdAt := time.Unix(int64(lastFMuser.User.Registered.Text), 0).Format(time.RFC850)

	//m := fmt.Sprintf("<b>%s</b>\n\n", lastFMuser.User.Name)
	m := mdparser.GetBold(lastFMuser.User.Name).AppendNormal("\n\n")

	//m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)
	m = m.AppendBold("Playcount: ").AppendNormal(lastFMuser.User.Playcount + "\n")

	//m += fmt.Sprintf("<b>Playlist Count</b>: %s\n", lastFMuser.User.Playlists)
	m = m.AppendBold("Playlist Count: ").AppendNormal(lastFMuser.User.Playlists + "\n")

	m = m.AppendBold("Gender: ")
	if len(lastFMuser.User.Gender) < 1 {
		m = m.AppendNormal("N/A\n")
	} else {
		m = m.AppendNormal(lastFMuser.User.Gender + "\n")
	}
	//m += fmt.Sprintf("<b>Gender</b>: %s\n", lastFMuser.User.Gender)

	m = m.AppendBold("Playcount: ").AppendNormal(lastFMuser.User.Playcount + "\n")

	// m += fmt.Sprintf("<b>Playcount</b>: %s\n", lastFMuser.User.Playcount)

	if lastFMuser.User.Age != "0" {
		m = m.AppendBold("Age: ").AppendNormal(lastFMuser.User.Age + "\n\n")
	}
	//m += fmt.Sprintf("<b>Age</b>: %s\n\n", lastFMuser.User.Age)

	m = m.AppendBold("Created at\n ").AppendMono(" " + createdAt + "\n")

	//m += fmt.Sprintf("<b>Created at</b>\n <code>%s</code>", createdAt)

	pic := getPicUrl(lastFMuser.User.Image)
	if pic == nil {
		//_, err = msg.Reply(b, m, &gotgbot.SendMessageOpts{ParseMode: "html"})
		_, err = msg.Reply(b, m.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		return err
	}

	_, err = b.SendPhoto(msg.Chat.Id, *pic, &gotgbot.SendPhotoOpts{ParseMode: "markdownv2", Caption: m.ToString(),
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "View on Last.FM", Url: lastFMuser.User.URL},
		}}}})
	return err
}

func getPicUrl(images []lastfm.Image) *string {
	for _, image := range images {
		if image.Size == lastfm.Extralarge {
			if len(image.Text) < 2 {
				return nil
			}
			return &image.Text
		}
	}

	return nil
}
