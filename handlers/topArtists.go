package handlers

import (
	"fmt"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"strconv"
)

func topArtistsHandler(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	msg := ctx.Message
	args := ctx.Args()
	var l int
	if len(args) == 1 {
		l = 5
		_, err := msg.Reply(b, "Amount not provided, sending first 5 results.", &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	} else {
		l, err = strconv.Atoi(args[1])
		if err != nil {
			_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()), &gotgbot.SendMessageOpts{ParseMode: "html"})
			if err != nil {
				return err
			}
		}
	}
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

	mtxt := mdparser.GetBold("Current top artists on Last.FM").AppendNormal("\n")
	//lnk, err := url.ParseQuery(fmt.Sprintf("tg://user?id=%d", msg.From.Id))
	mtxt = mtxt.AppendItalic("As requested by ").AppendMention(msg.From.FirstName, msg.From.Id).AppendNormal("\n\n")
	artists := d.Artists.Artist
	for i := 0; i < l; i++ {

		mtxt = mtxt.AppendHyperLink(artists[i].Name, artists[i].URL).AppendNormal("\n")
		mtxt = mtxt.AppendMono(artists[i].Playcount).AppendNormal(" ").AppendItalic("total plays.").AppendNormal("\n")
		mtxt = mtxt.AppendMono(artists[i].Listeners).AppendNormal(" ").AppendItalic("total listeners.").AppendNormal("\n\n")
		/*
			mtxt += fmt.Sprintf("[%s](%s)\n", artists[1].Name, artists[i].URL)
			mtxt += fmt.Sprintf("%s _Total Plays_\n", artists[i].Playcount)
			mtxt += fmt.Sprintf("%s _Total Listeners_\n\n", artists[i].Listeners)
		*/
	}
	_, err = b.SendMessage(msg.Chat.Id, mtxt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2", DisableWebPagePreview: true})
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	return nil

}
