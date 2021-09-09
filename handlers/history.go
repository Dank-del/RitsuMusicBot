package handlers

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
)

//var historyMap map[string]interface{}

func historyCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
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

	m := mdparser.GetBold("Recently played tracks by ").AppendMention(user.FirstName, user.Id).AppendNormal("\n\n")
	for a, e := range tracks {
		m = m.AppendNormal(fmt.Sprintf("%d", a+1)).AppendNormal(": ")
		m = m.AppendHyperLink(fmt.Sprintf("%s - %s", e.Artist.Name, e.Name), e.URL)
		if e.Loved == "1" {
			m = m.AppendItalic(" (Loved ♥)")
		}
		m = m.AppendNormal("\n")
		if e.Album.Text != "" {
			m = m.AppendBold("From album: ")
			m = m.AppendItalic(fmt.Sprintf("%s\n", e.Album.Text))
		}
		if a > 20 {
			break
		}
	}

	cb := "history_" + fuckAbs(ctx.EffectiveChat.Id) +
		"_" + fuckAbs(ctx.EffectiveUser.Id)

	kb := &gotgbot.InlineKeyboardMarkup{}
	kb.InlineKeyboard = append(kb.InlineKeyboard, []gotgbot.InlineKeyboardButton{})
	kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "\u23ea",
		CallbackData: cb + "_back",
	})
	kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "\u23e9",
		CallbackData: cb + "_next",
	})
	// fmt.Println(m)
	_, err = msg.Reply(b, m.ToString(),
		&gotgbot.SendMessageOpts{
			ParseMode:             "markdownv2",
			DisableWebPagePreview: true,
			ReplyMarkup:           kb,
		})
	return err
}

func fuckAbs(i int64) string {
	return strconv.FormatInt(int64(math.Abs(float64(i))), 10)
}

//  func(cq *gotgbot.CallbackQuery)
func historyCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, "history_")
}

// type Response func(b *gotgbot.Bot, ctx *ext.Context) error
func historyCallBackResponse(b *gotgbot.Bot, ctx *ext.Context) error {

	return nil
}
