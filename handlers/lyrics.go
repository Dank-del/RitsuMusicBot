package handlers

import (
	"fmt"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/google/uuid"
	"gitlab.com/Dank-del/lastfm-tgbot/core/utilities"
	"strings"

	"github.com/Dank-del/MusixScrape/musixScrape"
	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func lyricsInlineFilter(q *gotgbot.InlineQuery) bool {
	if q == nil {
		return false
	}
	return strings.Contains(strings.ToLower(q.Query), lyricsCommand)
}

func lyricsInline(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	inlq := ctx.InlineQuery
	q := strings.Trim(inlq.Query, lyricsCommand)
	var results []gotgbot.InlineQueryResult
	if q == "" {
		results = append(results, gotgbot.InlineQueryResultArticle{Id: inlq.Id, Title: "Query required",
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: "No query provided"}})
		_, err := inlq.Answer(b,
			results,
			&gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
		if err != nil {
			return err
		}
	}
	var l []musixScrape.SearchResult

	for len(l) < 1 {
		l, err = config2.Local.MusixMatchSession.Search(q)
		if err != nil {
			return err
		}
	}

	for i := range l {
		var link string
		keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
		txt := mdparser.GetBold("Results for " + q).Normal("\n\n")
		res, err := config2.Local.MusixMatchSession.GetLyrics(l[i].Url)
		if err != nil {
			txt = txt.Italic(err.Error()).Normal("\n")
		} else {
			link, err = utilities.PostLyrics(res.Song, res.Artist, res.Lyrics, config2.Local.TelegraphClient)
			if err != nil || link == strongStringGo.EMPTY {
				txt = mdparser.GetItalic("Failed to send lyrics due to telegraph error / Lyrics being empty")
				keyboard = nil
			} else {
				keyboard[0] = append(keyboard[0], gotgbot.InlineKeyboardButton{Text: "Telegra.ph", Url: link})
			}
		}
		var markup *gotgbot.InlineKeyboardMarkup
		if keyboard != nil {
			markup = &gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}
		}
		txt = mdparser.GetBold(res.Artist).Normal(" - ").Italic(res.Song).Bold(" [Lyrics]")
		results = append(results, gotgbot.InlineQueryResultArticle{
			Id:    uuid.NewString(),
			Title: fmt.Sprintf("%s - %s", l[i].Artist, l[i].Song),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: txt.ToString(),
				ParseMode:   "markdownv2",
			},
			ReplyMarkup: markup,
		})
	}

	_, err = inlq.Answer(b,
		results,
		&gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	if err != nil {
		return err
	}
	return nil

}

func lyricsHandler(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	msg := ctx.EffectiveMessage
	args := ctx.Args()
	_, err = b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	if len(args) == 1 {
		_, err := msg.Reply(b, mdparser.GetItalic("Query required").ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		return err
	}

	q := strings.Join(args[1:], " ")
	// fmt.Println(q)
	// var l []musixScrape.LyricResult
	// e := 0
	l, err := config2.Local.MusixMatchSession.Search(q)
	if err != nil {
		errm := mdparser.GetBold("Failed due to: ").Italic(err.Error())
		_, err := msg.Reply(b, errm.ToString(), config2.GetDefaultMdOpt())
		return err
	} else if len(l) == 0 {
		errm := mdparser.GetItalic("No results found")
		_, err := msg.Reply(b, errm.ToString(), config2.GetDefaultMdOpt())
		return err
	}
	txt := mdparser.GetItalic(l[0].Artist).Normal(" - ").AppendBoldThis(l[0].Song).Normal("\n")
	res, err := config2.Local.MusixMatchSession.GetLyrics(l[0].Url)
	link := strongStringGo.EMPTY
	if err != nil {
		txt.Normal(err.Error())
	} else {
		link, err = utilities.PostLyrics(res.Song, res.Artist, res.Lyrics, config2.Local.TelegraphClient)
		if err != nil || link == strongStringGo.EMPTY {
			_, err = msg.Reply(b,
				mdparser.GetItalic("Failed to send lyrics due to telegraph error / Lyrics being empty").ToString(),
				config2.GetDefaultMdOpt())
			return err
		}
	}
	keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
	keyboard[0] = append(keyboard[0], gotgbot.InlineKeyboardButton{Text: "Telegra.ph", Url: link})
	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2",
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}})
	return err
}
