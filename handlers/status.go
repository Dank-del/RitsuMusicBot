package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"html"
	"net/url"
	"strings"
)

func statusFilter(msg *gotgbot.Message) bool {
	return strings.Contains(strings.ToLower(msg.Text), "status")
}

func statusHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	uname, err := database.GetUser(msg.From.Id)
	_, err = b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}

	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", html.EscapeString(err.Error())), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername)
	if err != nil {
		logging.Warn(err.Error())
		return err
	}

	if d.Error != 0 {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", d.Message), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}

	if d.Recenttracks == nil {
		_, err := msg.Reply(b, "<i>You haven't scrobbed anything yet...</i>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	var s string
	if d.Recenttracks.Track[0].Attr != nil {
		s = "is now"
	} else {
		s = "was"
	}
	track := d.Recenttracks.Track[0]
	lfmUser, err := lastfm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.Warn(err.Error())
		return err
	}
	m := fmt.Sprintf("%s %s listening to\n", html.EscapeString(msg.From.FirstName), s)
	m += fmt.Sprintf("<i>%s</i> - <b>%s\n</b>", html.EscapeString(track.Artist.Text), track.Name)
	m += fmt.Sprintf("<i>%s total plays</i>", lfmUser.User.Playcount)
    yturl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s",
    	url.QueryEscape(fmt.Sprintf("%s - %s", track.Artist.Text, track.Name)))
	_, err = msg.Reply(b, m,
		&gotgbot.SendMessageOpts{ParseMode: "html", ReplyMarkup:
			gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "View on Last.FM", Url: track.URL},
				{Text: "Youtube", Url: yturl},
			}}}})
	return err
}
