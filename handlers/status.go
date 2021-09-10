package handlers

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	genius "gitlab.com/Dank-del/lastfm-tgbot/lyrics"
)

func statusFilter(msg *gotgbot.Message) bool {
	return strings.HasPrefix(strings.ToLower(msg.Text), statusMessage)
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

	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
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

	md := mdparser.GetHyperLink(msg.From.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n") //mdparser.GetNormal(fmt.Sprintf("%s %s listening to", msg.From.FirstName, s)).AppendNormal("\n")
	md = md.AppendItalic(track.Artist.Name).AppendNormal(" - ").AppendBold(track.Name).AppendNormal("\n")
	md = md.AppendItalic(fmt.Sprintf("%s total plays", lfmUser.User.Playcount))
	if track.Loved == "1" {
		md = md.AppendNormal(", ").AppendItalic("Loved ♥")
	}
	m := md.ToString()
	/*
		m := fmt.Sprintf("%s %s listening to\n", html.EscapeString(msg.From.FirstName), s)
		m += fmt.Sprintf("<i>%s</i> - <b>%s\n</b>", html.EscapeString(track.Artist.Name), track.Name)
		m += fmt.Sprintf("<i>%s total plays</i>", lfmUser.User.Playcount)
		if track.Loved == "1" {
			m += fmt.Sprintf(", <i>Loved ♥</i>")
		}
	*/
	yturl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s",
		url.QueryEscape(fmt.Sprintf("%s - %s", track.Artist.Name, track.Name)))
	_, err = msg.Reply(b, m,
		&gotgbot.SendMessageOpts{ParseMode: "markdownv2", ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "View on Last.FM", Url: track.URL},
			{Text: "Youtube", Url: yturl},
		}}}, DisableWebPagePreview: true})

	if strings.Contains(msg.Text, lyricsCommand) {
		// m = fmt.Sprintf("<b>Lyrics: %s - %s</b>\n\n", html.EscapeString(track.Artist.Name), html.EscapeString(track.Name))
		m := mdparser.GetBold(fmt.Sprintf("Lyrics: %s - %s", track.Artist.Name, track.Name)).AppendNormal("\n\n")
		var l []string
		e := 0
		for len(l) < 5 {
			l, err = genius.GetLyrics(fmt.Sprintf("%s - %s", track.Artist.Name, track.Name))
			e++
			if e > 10 {
				_, err := b.SendMessage(msg.Chat.Id, mdparser.GetItalic(err.Error()).ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
				return err
			}
		}
		if err != nil {
			// m += fmt.Sprintf("<i>Error: %s</i>", err.Error())
			m = m.AppendItalic(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			for i := range l {
				// m += fmt.Sprintf("<i>%s</i>\n", html.EscapeString(l[i]))
				m = m.AppendItalic(l[i]).AppendNormal("\n")
			}
		}
		_, err = msg.Reply(b, m.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	}
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
	// fmt.Println(user)
	var m mdparser.WMarkDown
	// fmt.Println(m)
	var results []gotgbot.InlineQueryResult

	uname, err := database.GetUser(user.Id)
	// fmt.Println(uname)
	if err != nil || uname.LastFmUsername == "" {
		m = mdparser.GetBold(user.FirstName).AppendNormal(" ").AppendItalic("haven't registered themselves on this bot yet.").AppendNormal("\n")
		m = m.AppendBold("Please use ").AppendMono("/setusername")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s's needs to register themselves", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
	if err != nil {
		logging.Warn(err.Error())
		return err
	}

	if d.Error != 0 {
		m = mdparser.GetItalic(fmt.Sprintf("Error: %s", d.Message))
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("Sorry %s!, I encountered an error :/", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	if d.Recenttracks == nil {
		m = mdparser.GetItalic("No tracks were being played.")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s haven't played anything yet.", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	lfmUser, err := lastfm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.Warn(err.Error())
		return err
	}

	for e, i := range d.Recenttracks.Track {
		var s string
		if d.Recenttracks.Track[e].Attr != nil {
			s = "is now"
		} else {
			s = "was"
		}
		m := mdparser.GetHyperLink(user.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n") //mdparser.GetNormal(fmt.Sprintf("%s %s listening to", user.FirstName, s)).AppendNormal("\n")
		m = m.AppendItalic(i.Artist.Name).AppendNormal(" - ").AppendBold(i.Name).AppendNormal("\n")
		m = m.AppendItalic(fmt.Sprintf("%s total plays", lfmUser.User.Playcount))
		if i.Loved == "1" {
			m = m.AppendNormal(", ").AppendItalic("Loved ♥")
		}

		/*if strings.Contains(query.Query, "lyrics") {
			m = m.AppendNormal("\n\n").AppendBold("Lyrics").AppendNormal("\n")
			l, err := LyricsClient.Search(i.Artist.Name, i.Name)
			if err != nil {
				m = m.AppendItalic(fmt.Sprintf("Error: %s", err.Error()))
			}
			m = m.AppendItalic(l)
		} */
		results = append(results, gotgbot.InlineQueryResultArticle{Id: uuid.New().String(), Title: fmt.Sprintf("%s - %s", i.Artist.Name, i.Name),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2", DisableWebPagePreview: true}})

		if e > 12 {
			break
		}
	}

	if err != nil {
		_, err := query.Answer(b,
			results,
			&gotgbot.AnswerInlineQueryOpts{})
		if err != nil {
			return err
		}
	}
	_, err = query.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	return nil

}

func getStatus(user *gotgbot.User) (string, error) {
	uname, err := database.GetUser(user.Id)
	if uname.LastFmUsername == "" {
		m := mdparser.GetBold(user.FirstName).AppendNormal(" ").AppendItalic("haven't registered themselves on this bot yet.").AppendNormal("\n")
		m = m.AppendBold("Please use ").AppendMono("/setusername")
		return m.ToString(), err
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
	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
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
			return mdparser.GetItalic("You haven't scrobbed anything yet...").ToString(), err
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
	m := mdparser.GetHyperLink(user.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n") //mdparser.GetNormal(fmt.Sprintf("%s %s listening to", user.FirstName, s)).AppendNormal("\n")
	m = m.AppendItalic(track.Artist.Name).AppendNormal(" - ").AppendBold(track.Name).AppendNormal("\n")
	m = m.AppendItalic(fmt.Sprintf("%s total plays", lfmUser.User.Playcount))
	if track.Loved == "1" {
		m = m.AppendNormal(", ").AppendItalic("Loved ♥")
	}
	/*
		m := fmt.Sprintf("%s %s listening to\n", html.EscapeString(user.FirstName), s)
		m += fmt.Sprintf("<i>%s</i> - <b>%s\n</b>", html.EscapeString(track.Artist.Name), track.Name)
		m += fmt.Sprintf("<i>%s total plays</i>", lfmUser.User.Playcount)
		if track.Loved == "1" {
			m += fmt.Sprintf(", <i>Loved ♥</i>")
		}*/
	return m.ToString(), err
}
