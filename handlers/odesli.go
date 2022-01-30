package handlers

import (
	"regexp"

	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/libs/odesli"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func msgLinkFilter(msg *gotgbot.Message) bool {
	txt := msg.Text
	m, err := regexp.MatchString(urlRegEx, txt)

	return err != nil && m
}

func odesliMessageHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	re, err := regexp.Compile(`(?:(?:https?|ftp):\/\/)?[\w/\-?=%.]+\.[\w/\-&?=%.]+`)
	if err != nil {
		return ext.ContinueGroups
	}

	urls := re.FindAllString(msg.Text, -1)
	if urls == nil {
		return nil
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

	_, err = b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}

	txt := mdparser.GetUserMention(msg.From.FirstName, msg.From.Id).Bold(" sent").Normal("\n\n")
	t := d.EntitiesByUniqueID[d.EntityUniqueID]
	txt.Italic(t.ArtistName).Normal(" - ").Bold(t.Title).Normal("\n\n")
	links := d.LinksByPlatform
	dot := "• "
	if links.Deezer != nil {
		txt.Normal(dot).Link("Deezer", links.Deezer.URL).ElThis()
	}

	if links.Itunes != nil {
		txt.Normal(dot).Link("Itunes", links.Itunes.URL).ElThis()
	}

	if links.Tidal != nil {
		txt.Normal(dot).Link("Tidal", links.Deezer.URL).ElThis()
	}

	if links.AmazonMusic != nil {
		txt.Normal(dot).Link("Amazon Music", links.AmazonMusic.URL).ElThis()
	}

	if links.AmazonStore != nil {
		txt.Normal(dot).Link("Amazon Store", links.AmazonStore.URL).ElThis()
	}

	if links.Napster != nil {
		txt.Normal(dot).Link("Napster", links.Napster.URL).ElThis()
	}
	if links.AppleMusic != nil {
		txt.Normal(dot).Link("Apple Music", links.AppleMusic.URL).ElThis()
	}
	if links.Pandora != nil {
		txt.Normal(dot).Link("Pandora", links.Pandora.URL).ElThis()
	}
	if links.Youtube != nil {
		txt.Normal(dot).Link("Youtube", links.Youtube.URL).ElThis()
	}
	if links.Spotify != nil {
		txt.Normal(dot).Link("Spotify", links.Spotify.URL).ElThis()
	}

	if links.Soundcloud != nil {
		txt.Normal(dot).Link("SoundCloud", links.Soundcloud.URL).ElThis()
	}

	if links.Yandex != nil {
		txt.Normal(dot).Link("Yandex", links.Yandex.URL).ElThis()
	}

	if links.YoutubeMusic != nil {
		txt.Normal(dot).Link("Youtube Music", links.YoutubeMusic.URL).ElThis()
	}

	txt.Normal("\n\n").Italic("Powered by Odesli smartlink API")

	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             "markdownv2",
		DisableWebPagePreview: true,
	})
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}

	return ext.EndGroups
}
