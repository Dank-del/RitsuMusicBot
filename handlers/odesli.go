package handlers

import (
	"errors"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/odesli"
	"regexp"
)

func msgLinkFilter(msg *gotgbot.Message) bool {
	txt := msg.Text
	m, err := regexp.MatchString(urlRegEx, txt)
	if err != nil {
		return false
	}
	return m
}

func odesliMessageHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	d, err := odesli.GetLinks(msg.Text)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	if d.Code != "" {
		logging.Error(d.Code)
		return errors.New(d.Code)
	}
	txt := mdparser.GetUserMention(msg.From.FirstName, msg.From.Id).AppendBold(" sent").AppendNormal("\n")
	t := d.EntitiesByUniqueID[d.EntityUniqueID]
	txt = txt.AppendItalic(t.ArtistName).AppendNormal(" - ").AppendBold(t.Title).AppendNormal("\n")
	links := d.LinksByPlatform
	if links.Deezer != nil {
		txt = txt.AppendHyperLink("Deezer", links.Deezer.URL).AppendNormal(" ")
	}
	if links.Itunes != nil {
		txt = txt.AppendHyperLink("Itunes", links.Itunes.URL).AppendNormal(" ")
	}
	if links.Tidal != nil {
		txt = txt.AppendHyperLink("Tidal", links.Deezer.URL).AppendNormal(" ")
	}
	if links.AmazonMusic != nil {
		txt = txt.AppendHyperLink("Amazon Music", links.AmazonMusic.URL).AppendNormal(" ")
	}
	if links.AmazonStore != nil {
		txt = txt.AppendHyperLink("Amazon Store", links.AmazonStore.URL).AppendNormal(" ")
	}
	if links.Napster != nil {
		txt = txt.AppendHyperLink("Napster", links.Napster.URL).AppendNormal(" ")
	}
	if links.AppleMusic != nil {
		txt = txt.AppendHyperLink("Apple Music", links.AppleMusic.URL).AppendNormal(" ")
	}
	if links.Pandora != nil {
		txt = txt.AppendHyperLink("Pandora", links.Pandora.URL).AppendNormal(" ")
	}
	if links.Youtube != nil {
		txt = txt.AppendHyperLink("Youtube", links.Youtube.URL).AppendNormal(" ")
	}
	if links.Spotify != nil {
		txt = txt.AppendHyperLink("Spotify", links.Spotify.URL).AppendNormal(" ")
	}
	if links.Soundcloud != nil {
		txt = txt.AppendHyperLink("SoundCloud", links.Soundcloud.URL).AppendNormal(" ")
	}
	if links.Yandex != nil {
		txt = txt.AppendHyperLink("Yandex", links.Yandex.URL).AppendNormal(" ")
	}
	if links.YoutubeMusic != nil {
		txt = txt.AppendHyperLink("Youtube Music", links.YoutubeMusic.URL).AppendNormal(" ")
	}
	// print(txt.ToString())
	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	return nil
}
