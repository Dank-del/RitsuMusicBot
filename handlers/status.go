package handlers

import (
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"

	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/libs/last.fm"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
)

func statusFilter(msg *gotgbot.Message) bool {
	chatID := msg.Chat.Id
	d, err := database.GetChat(chatID)
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return false
	}
	// print(d.GetStatusMessage())
	return strings.HasPrefix(strings.ToLower(msg.Text), d.GetStatusMessage())
}

func statusHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	uname, err := database.GetLastFMUserFromDB(msg.From.Id)
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

		return nil
	}

	d, err := last_fm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	if d.Error != 0 {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", d.Message), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
		return nil
	}

	if d.Recenttracks == nil || len(d.Recenttracks.Track) == 0 {
		_, err := msg.Reply(b, "<i>You haven't scrobbed anything yet...</i>",
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}

		return nil
	}
	var s string
	if d.Recenttracks.Track[0].Attr != nil {
		s = "is now"
	} else {
		s = "was"
	}

	track := &d.Recenttracks.Track[0]
	lfmUser, err := last_fm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}
	setting, err := database.GetBotUserByID(msg.From.Id)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	var md mdparser.WMarkDown
	var pic *string
	if track.Image != nil {
		pic = getPicUrl(track.Image)
	}

	if pic == nil && track.Album != nil {
		pic = getPicUrl(track.Album.Image)
	}

	if pic == nil && track.Artist != nil {
		pic = getPicUrl(track.Artist.Image)
	}

	hasAlbum := track.Album != nil && pic != nil
	md = mdparser.GetNormal("🎧 ")
	if hasAlbum {
		md = md.AppendHyperLink("\u2063", *pic)
	}

	if setting.ShowProfile {
		md = md.AppendHyperLink(msg.From.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
	} else {
		md = md.AppendBold(msg.From.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
	}

	md = md.AppendItalic(track.Artist.Name).AppendNormal(" - ").AppendBold(track.Name).AppendNormal("\n")
	topTracks, err := last_fm.GetTopTracks(uname.LastFmUsername)
	if err != nil {
		md.AppendItalic(err.Error()).AppendNormal("\n")
	}
	if topTracks.Error != 0 {
		md.AppendItalic("Error fetching recent tracks: " + topTracks.Message)
	} else {
		var scb string
		for _, e := range topTracks.Toptracks.Track {
			if strings.Compare(track.Name, e.Name) == 0 && strings.Compare(track.Artist.Name, e.Artist.Name) == 0 {
				scb = e.Playcount
				break
			}
		}
		if scb != "" {
			md.AppendItalic("Scrobbled " + scb + " times by " + msg.From.FirstName).AppendNormal("\n")
		}
	}
	if track.Loved == "1" {
		md = md.AppendItalic("Loved ♥").AppendNormal("\n")
	}
	md = md.AppendItalic(fmt.Sprintf("%s scrobbles so far", lfmUser.User.Playcount))
	topTracks = nil

	_, err = msg.Reply(b, md.ToString(),
		&gotgbot.SendMessageOpts{
			ParseMode:             "markdownv2",
			DisableWebPagePreview: true,
			ReplyMarkup:           generateButtons(track, hasAlbum, msg.From.Id),
		})

	if strings.Contains(msg.Text, lyricsCommand) {
		m := mdparser.GetBold(fmt.Sprintf("Lyrics: %s - %s", track.Artist.Name, track.Name)).AppendNormal("\n\n")
		// var l []musixScrape.LyricResult
		l, err := config2.Local.MusixMatchSession.Search(fmt.Sprintf("%s - %s", track.Artist.Name, track.Name))
		if err != nil {
			errm := mdparser.GetBold("Failed due to: ").AppendItalic(err.Error())
			_, err := msg.Reply(b, errm.ToString(), config2.GetDefaultMdOpt())
			return err
		} else if len(l) == 0 {
			errm := mdparser.GetItalic("No results found")
			_, err := msg.Reply(b, errm.ToString(), config2.GetDefaultMdOpt())
			return err
		}
		m.Italic(l[0].Artist).Normal(" - ").AppendBoldThis(l[0].Song).Normal("\n")
		m.Normal(l[0].Lyrics)
		_, err = msg.Reply(b, m.ToString(), config2.GetDefaultMdOpt())
		return err
	}
	return err
}

func generateButtons(track *last_fm.Track, album bool,
	id int64) *gotgbot.InlineKeyboardMarkup {
	yturl := fmt.Sprintf("https://www.youtube.com/results?search_query=%s",
		url.QueryEscape(fmt.Sprintf("%s - %s", track.Artist.Name, track.Name)))
	var tmpmarkup gotgbot.InlineKeyboardButton
	keyboard := make([][]gotgbot.InlineKeyboardButton, 2)

	// view on "Last .FM" button.
	tmpmarkup = gotgbot.InlineKeyboardButton{
		Text: "View on Last.FM",
		Url:  track.URL,
	}
	keyboard[0] = append(keyboard[0], tmpmarkup)

	tmpmarkup = gotgbot.InlineKeyboardButton{
		Text: "Youtube",
		Url:  yturl,
	}
	keyboard[0] = append(keyboard[0], tmpmarkup)

	if album {
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text:         albumText,
			CallbackData: albumPrefix + strconv.FormatInt(id, 10),
		}
		keyboard[1] = append(keyboard[1], tmpmarkup)
	}

	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

//  func(cq *gotgbot.CallbackQuery)
func albumCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, albumPrefix)
}

// type Response func(b *gotgbot.Bot, ctx *ext.Context) error
func albumCallBackResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	mystrs := strings.Split(ctx.CallbackQuery.Data, "_")
	id, err := strconv.ParseInt(mystrs[1], 10, 64)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	if id != ctx.EffectiveUser.Id {
		_, err = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "this is not for you",
			ShowAlert: true,
		})

		if err != nil {
			logging.SUGARED.Error(err.Error())
		}
		return err
	}

	msg := ctx.EffectiveMessage
	preview := msg.ReplyMarkup.InlineKeyboard[1][0].Text == albumText
	if preview {
		msg.ReplyMarkup.InlineKeyboard[1][0].Text = hideText
	} else {
		msg.ReplyMarkup.InlineKeyboard[1][0].Text = albumText
	}

	_, _, err = msg.EditText(b, msg.Text, &gotgbot.EditMessageTextOpts{
		Entities:              ctx.EffectiveMessage.Entities,
		DisableWebPagePreview: !preview,
		ReplyMarkup:           *msg.ReplyMarkup,
	})

	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}

	_, err = b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: "Toggled."})
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

	uname, err := database.GetLastFMUserFromDB(user.Id)
	// fmt.Println(uname)
	if err != nil || !uname.IsValid() {
		m = mdparser.GetBold(user.FirstName).AppendNormal(" ").AppendItalic("haven't registered themselves on this bot yet.").AppendNormal("\n")
		m = m.AppendBold("Please use ").AppendMono("/setusername")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s's needs to register themselves", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
		_, _ = query.Answer(b, results, nil)
		return ext.EndGroups
	}

	d, err := last_fm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	if d.Error != 0 {
		m = mdparser.GetItalic(fmt.Sprintf("Error: %s", d.Message))
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("Sorry %s!, I encountered an error :/", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	if d.Recenttracks == nil || len(d.Recenttracks.Track) == 0 {
		m = mdparser.GetItalic("No tracks were being played.")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s haven't played anything yet.", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	lfmUser, err := last_fm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	for e, i := range d.Recenttracks.Track {
		var s string
		if d.Recenttracks.Track[e].Attr != nil {
			s = "is now"
		} else {
			s = "was"
		}
		setting, err := database.GetBotUserByID(user.Id)
		if err != nil {
			logging.SUGARED.Warn(err.Error())
			return err
		}

		var m mdparser.WMarkDown
		switch setting.ShowProfile {
		case true:
			m = mdparser.GetHyperLink(user.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		case false:
			m = mdparser.GetBold(user.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		default:
			m = mdparser.GetBold(user.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		} //mdparser.GetNormal(fmt.Sprintf("%s %s listening to", user.FirstName, s)).AppendNormal("\n")
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
		_, _ = query.Answer(b, results, nil)
		return ext.EndGroups
	}
	_, err = query.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}
	return nil

}

func getStatus(user *gotgbot.User) (mdparser.WMarkDown, error) {
	uname, err := database.GetLastFMUserFromDB(user.Id)
	if uname.LastFmUsername == "" {
		m := mdparser.GetBold(user.FirstName).AppendNormal(" ").AppendItalic("haven't registered themselves on this bot yet.").AppendNormal("\n")
		m = m.AppendBold("Please use ").AppendMono("/setusername")
		return m, err
	}
	// fmt.Println(err)
	if err != nil {
		return nil, err
	}

	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	d, err := last_fm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return nil, err
	}

	if d.Error != 0 {
		if err != nil {
			logging.SUGARED.Warn(err.Error())
			return nil, err
		}
	}

	if d.Recenttracks == nil || len(d.Recenttracks.Track) == 0 {
		if err != nil {
			return mdparser.GetItalic("You haven't scrobbed anything yet..."), err
		}
	}

	var s string
	if d.Recenttracks.Track[0].Attr != nil {
		s = "is now"
	} else {
		s = "was"
	}
	track := d.Recenttracks.Track[0]
	lfmUser, err := last_fm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return nil, err
	}

	setting, err := database.GetBotUserByID(user.Id)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return nil, err
	}
	var m mdparser.WMarkDown
	if setting.ShowProfile {
		m = mdparser.GetHyperLink(user.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
	} else {
		m = mdparser.GetBold(user.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
	}

	m = m.AppendItalic(track.Artist.Name).AppendNormal(" - ").AppendBold(track.Name).AppendNormal("\n")
	m = m.AppendItalic(fmt.Sprintf("%s total plays", lfmUser.User.Playcount))
	if track.Loved == "1" {
		m = m.AppendNormal(", ").AppendItalic("Loved ♥")
	}

	return m, err
}
