package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"strings"
)

func topArtistsHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	d, err := lastfm.GetTopArtists()
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}

	mtxt := "*Current top artists on Last.FM*\n\n"
	//lnk, err := url.ParseQuery(fmt.Sprintf("tg://user?id=%d", msg.From.Id))
	mtxt += fmt.Sprintf("_As requested by_  [%s](tg://user?id=%d)\n", msg.From.FirstName, msg.From.Id)
	artists := d.Artists.Artist
	for i := 0; i < 5; i++ {
		mtxt += fmt.Sprintf("[%s](%s)\n", artists[1].Name, artists[i].URL)
		mtxt += fmt.Sprintf("%s _Total Plays_\n", artists[i].Playcount)
		mtxt += fmt.Sprintf("%s _Total Listeners_\n\n", artists[i].Listeners)
	}
	fmt.Println(strings.Trim(mtxt, ""))
	a, err := b.SendMessage(msg.Chat.Id, strings.Trim(mtxt, ""), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	fmt.Println(a)
	if err != nil {
		return err
	}
	return nil

}
