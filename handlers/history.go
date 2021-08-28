package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"html"
)

func historyCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From

	dbuser, err := database.GetUser(user.Id)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}
	if dbuser.LastFmUsername == "" {
		_, err := msg.Reply(b, "<i>You haven't registered yourself on this bot yet</i>\n<b>Use /setusername</b>",
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}
	grc, err := last_fm.GetRecentTracksByUsername(dbuser.LastFmUsername)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	if grc.Error != 0 {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", grc.Message),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	tracks := grc.Recenttracks.Track

	m := fmt.Sprintf("<b> %s's recent tracks</b>\n\n", html.EscapeString(user.FirstName))
	for a, e := range tracks {
		m += fmt.Sprintf("%d - <a href=%s>%s</a>\n", a+1, e.URL,
			html.EscapeString(fmt.Sprintf("%s - %s", e.Artist.Text, e.Name)))
		m += fmt.Sprintf("<i>From %s</i>\n", html.EscapeString(e.Album.Text))
		if a > 7 {
			break
		}
	}
	fmt.Println(m)
	q, err := msg.Reply(b, m,
		&gotgbot.SendMessageOpts{ParseMode: "html"})
	fmt.Println(q)
	return err
}
