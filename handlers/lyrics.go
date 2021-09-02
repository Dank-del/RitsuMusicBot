package handlers

import (
	"fmt"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	genius "gitlab.com/Dank-del/lastfm-tgbot/lyrics"
	"strings"
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
	var l []string
	e := 0
	txt := mdparser.GetBold(fmt.Sprintf("Results for %s", q)).AppendNormal("\n\n")
	for len(l) < 5 {
		l, err = genius.GetLyrics(q)
		if err != nil {
			return err
		}
		e++
		if e > 10 {
			results = append(results, gotgbot.InlineQueryResultArticle{Id: inlq.Id, Title: "Error",
				InputMessageContent: gotgbot.InputTextMessageContent{MessageText: mdparser.GetItalic(fmt.Sprintf("Error: %s", err.Error())).ToString(),
					ParseMode: "markdownv2"}})
			_, err := inlq.Answer(b,
				results,
				&gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
			if err != nil {
				return err
			}
		}
	}

	for i := range l {
		txt = txt.AppendItalic(l[i]).AppendNormal("\n")
	}

	results = append(results, gotgbot.InlineQueryResultArticle{Id: inlq.Id, Title: fmt.Sprintf("Lyrics: %s", q),
		InputMessageContent: gotgbot.InputTextMessageContent{MessageText: txt.ToString(),
			ParseMode: "markdownv2"}})
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
		_, err := b.SendMessage(msg.Chat.Id, mdparser.GetItalic("Query required").ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		return err
	}

	q := strings.Join(args[1:], " ")
	// fmt.Println(q)
	var l []string
	e := 0
	txt := mdparser.GetBold(fmt.Sprintf("Results for %s", q)).AppendNormal("\n\n")
	for len(l) < 5 {
		l, err = genius.GetLyrics(q)
		if err != nil {
			return err
		}
		e++
		if e > 10 {
			_, err := b.SendMessage(msg.Chat.Id, mdparser.GetItalic(err.Error()).ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
			return err
		}
	}

	for i := range l {
		txt = txt.AppendItalic(l[i]).AppendNormal("\n")
	}

	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	return err
}
