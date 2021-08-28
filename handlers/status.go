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
	return strings.Contains(strings.ToLower(msg.Text), statusMessage)
}

func statusHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	uname, err := database.GetUser(msg.From.Id)
	if uname.LastFmUsername == "" {
		_, err := msg.Reply(b, "<i>You haven't registered yourself on this bot yet</i>\n<b>Use /setusername</b>",
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
		return nil
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
	m += fmt.Sprintf("<i>%s</i> - <b>%s\n</b>", html.EscapeString(track.Artist.Name), track.Name)
	m += fmt.Sprintf("<i>%s total plays</i>", lfmUser.User.Playcount)
	if track.Loved == "1" {
		m += fmt.Sprintf(", <i>Loved ♥</i>")
	}
	yturl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s",
		url.QueryEscape(fmt.Sprintf("%s - %s", track.Artist.Name, track.Name)))
	_, err = msg.Reply(b, m,
		&gotgbot.SendMessageOpts{ParseMode: "html", ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "View on Last.FM", Url: track.URL},
			{Text: "Youtube", Url: yturl},
		}}}})
	return err
}

func statusInlineFilter(q *gotgbot.InlineQuery) bool {
	if q == nil {
		return false
	}
	return strings.Contains(strings.ToLower(q.Query), statusMessage)
}

func statusInline(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery
	user := query.From
	m, err := getStatus(&user)
	// fmt.Println(m)
	var results []gotgbot.InlineQueryResult
	results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s's status", user.FirstName),
		InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m, ParseMode: "html"}})
	if err != nil {
		_, err := query.Answer(b,
			results,
			&gotgbot.AnswerInlineQueryOpts{})
		if err != nil {
			return err
		}
	}
	_, err = query.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{CacheTime: 0})
	if err != nil {
		return err
	}
	return nil

}

func getStatus(user *gotgbot.User) (string, error) {
	uname, err := database.GetUser(user.Id)
	if uname.LastFmUsername == "" {
		return "<i>You haven't registered yourself on this bot yet</i>\n<b>Use /setusername</b>", err
	}
	// fmt.Println(err)
	if err != nil {
		return "", err
	}

	if err != nil {
		if err != nil {
			return "", err
		}
	}
	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername)
	if err != nil {
		logging.Warn(err.Error())
		return "", err
	}

	if d.Error != 0 {
		if err != nil {
			return d.Message, err
		}
	}

	if d.Recenttracks == nil {
		if err != nil {
			return "<i>You haven't scrobbed anything yet...</i>", err
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
		return "", err
	}

	m := fmt.Sprintf("%s %s listening to\n", html.EscapeString(user.FirstName), s)
	m += fmt.Sprintf("<i>%s</i> - <b>%s\n</b>", html.EscapeString(track.Artist.Name), track.Name)
	m += fmt.Sprintf("<i>%s total plays</i>", lfmUser.User.Playcount)
	if track.Loved == "1" {
		m += fmt.Sprintf(", <i>Loved ♥</i>")
	}
	return m, err
}
