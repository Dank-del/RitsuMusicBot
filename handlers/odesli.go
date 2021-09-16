package handlers

import (
	"regexp"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/odesli"
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
	re, err := regexp.Compile("(?:(?:https?|ftp):\\/\\/)?[\\w/\\-?=%.]+\\.[\\w/\\-&?=%.]+")
	if err != nil {
		return err
	}
	urls := re.FindAllString(msg.Text, -1)
	if urls == nil {
		return nil
	}
	_, err = b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	d, err := odesli.GetLinks(urls[0])
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}

	if d.Code != "" {
		// check if the link is valid or not. if the server returns
		// this error code, it means that the link is not a valid
		// link at all, so don't panic.
		if d.Code == "could_not_resolve_entity" {
			return ext.EndGroups
		}

		logging.SUGARED.Error(d.Code)
		return nil
	}
	txt := mdparser.GetUserMention(msg.From.FirstName, msg.From.Id).AppendBold(" sent").AppendNormal("\n\n")
	t := d.EntitiesByUniqueID[d.EntityUniqueID]
	txt = txt.AppendItalic(t.ArtistName).AppendNormal(" - ").AppendBold(t.Title).AppendNormal("\n\n")
	links := d.LinksByPlatform
	dot := "• "
	newline := "\n"
	if links.Deezer != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Deezer", links.Deezer.URL).AppendNormal(newline)
	}
	if links.Itunes != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Itunes", links.Itunes.URL).AppendNormal(newline)
	}
	if links.Tidal != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Tidal", links.Deezer.URL).AppendNormal(newline)
	}
	if links.AmazonMusic != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Amazon Music", links.AmazonMusic.URL).AppendNormal(newline)
	}
	if links.AmazonStore != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Amazon Store", links.AmazonStore.URL).AppendNormal(newline)
	}
	if links.Napster != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Napster", links.Napster.URL).AppendNormal(newline)
	}
	if links.AppleMusic != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Apple Music", links.AppleMusic.URL).AppendNormal(newline)
	}
	if links.Pandora != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Pandora", links.Pandora.URL).AppendNormal(newline)
	}
	if links.Youtube != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Youtube", links.Youtube.URL).AppendNormal(newline)
	}
	if links.Spotify != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Spotify", links.Spotify.URL).AppendNormal(newline)
	}
	if links.Soundcloud != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("SoundCloud", links.Soundcloud.URL).AppendNormal(newline)
	}
	if links.Yandex != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Yandex", links.Yandex.URL).AppendNormal(newline)
	}
	if links.YoutubeMusic != nil {
		txt = txt.AppendNormal(dot).AppendHyperLink("Youtube Music", links.YoutubeMusic.URL).AppendNormal(newline)
	}

	txt = txt.AppendNormal(newline + newline).AppendItalic("Powered by Odesli smartlink API")
	// print(txt.ToString())
	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2", DisableWebPagePreview: true})
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}
	return nil
}
