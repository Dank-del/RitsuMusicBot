package handlers

import (
	"fmt"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
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

	m := mdparser.GetBold(fmt.Sprintf("%s's recent tracks\n\n", user.FirstName))
	for a, e := range tracks {
		m = m.AppendNormal(fmt.Sprintf("%d", a+1)).AppendNormal(": ")
		m = m.AppendHyperLink(fmt.Sprintf("%s - %s\n", e.Artist.Text, e.Name), e.URL)
		m = m.AppendItalic(fmt.Sprintf("From %s\n", e.Album.Text))
		if a > 20 {
			break
		}
	}
	// fmt.Println(m)
	_, err = msg.Reply(b, m.ToString(),
		&gotgbot.SendMessageOpts{ParseMode: "markdownv2", DisableWebPagePreview: true})
	return err
}
