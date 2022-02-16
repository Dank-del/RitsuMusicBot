package handlers

import (
	"context"
	"fmt"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/Dank-del/MusixScrape/musixScrape"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/utilities"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"strings"
)

func spotifyNow(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	_, err := b.SendChatAction(ctx.EffectiveChat.Id, "typing")
	if err != nil {
		return err
	}
	msg, err := ctx.EffectiveMessage.Reply(b, "Please wait...", nil)
	ctxBg := context.Background()
	spotifyUser, err := database.GetSpotifyUser(user.Id)
	if err != nil {
		_, err = msg.Delete(b)
		if err != nil {
			return err
		}
		return linkSpotifyHandler(b, ctx)
	}
	if spotifyUser != nil {
		currentlyPlaying, err := spotifyUser.PlayerCurrentlyPlaying(ctxBg, spotify.Limit(1))
		if err != nil {
			return err
		}
		// currentUser, _ := spotifyUser.CurrentUser(ctxBg)
		if currentlyPlaying != nil && currentlyPlaying.Playing {
			/*
				    marshal, err := json.MarshalIndent(currentlyPlaying, " ", " ")
					if err != nil {
						return err
					}
					print(marshal)
			*/
			txt := mdparser.GetUserMention(user.FirstName, user.Id).Normal(" is listening to ")
			var artists []string
			for _, a := range currentlyPlaying.Item.Artists {
				artists = append(artists, a.Name)
			}
			img := currentlyPlaying.Item.Album.Images[0].URL
			txt.Link(currentlyPlaying.Item.Name,
				currentlyPlaying.Item.ExternalURLs["spotify"]).
				Normal(" by ").Bold(strings.Join(artists, ", "))
			keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
			keyboard[0] = append(keyboard[0], gotgbot.InlineKeyboardButton{
				Text: "Album",
				Url:  currentlyPlaying.Item.Album.ExternalURLs["spotify"],
			})
			var showLyric bool
			l, err := config2.Local.MusixMatchSession.Search(fmt.Sprintf("%s - %s", artists, currentlyPlaying.Item.Name))
			if len(l) == 0 || err != nil {
				showLyric = false
			} else {
				showLyric = true
			}
			var lyricLink string
			var res musixScrape.Lyrics
			if showLyric {
				res, err = config2.Local.MusixMatchSession.GetLyrics(l[0].Url)
				if err != nil || &res == nil || res.Lyrics == strongStringGo.EMPTY {
					showLyric = false
				}
				lyricLink, err = utilities.PostLyrics(res.Song, res.Artist, res.Lyrics, config2.Local.TelegraphClient)
				if err != nil || lyricLink == strongStringGo.EMPTY {
					showLyric = false
				}
			}
			if showLyric {
				keyboard[0] = append(keyboard[0], gotgbot.InlineKeyboardButton{
					Text: "Lyrics",
					Url:  lyricLink,
				})
			}
			_, err = msg.Delete(b)
			if err != nil {
				return err
			}
			_, err = b.SendPhoto(ctx.EffectiveChat.Id, img, &gotgbot.SendPhotoOpts{
				ParseMode:   "markdownv2",
				Caption:     txt.ToString(),
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard},
			})
			return err
		} else {
			_, err = msg.Delete(b)
			if err != nil {
				return err
			}
			_, err = ctx.EffectiveMessage.Reply(b, "You're not listening to anything..", nil)
			return err
		}
	}
	return ext.EndGroups
}

func spotifyInlineFilter(q *gotgbot.InlineQuery) bool {
	if q == nil {
		return false
	}
	return strings.Contains(strings.ToLower(q.Query), spotNowCommand)
}

func spotifyInline(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	inlq := ctx.InlineQuery
	ctxBq := context.Background()
	var results []gotgbot.InlineQueryResult
	var result gotgbot.InlineQueryResultArticle
	spotifyUser, err := database.GetSpotifyUser(inlq.From.Id)
	if err != nil {
		txt := mdparser.GetBold("Send me ").Mono("/" + linkSpotifyCommand).Bold(" to authenticate")
		result = gotgbot.InlineQueryResultArticle{
			Id:    uuid.NewString(),
			Title: "Not authorized",
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: txt.ToString(),
				ParseMode:   "markdownv2",
			},
		}
	} else {
		var txt mdparser.WMarkDown
		var title string
		var img string
		p, err := spotifyUser.PlayerCurrentlyPlaying(ctxBq)
		if err != nil {
			title = "Error"
			txt = mdparser.GetMono(err.Error())
		} else if !p.Playing {
			title = "Play something!"
			txt = mdparser.GetNormal("You're not listening to anything")
		} else {
			var artists []string
			img = p.Item.Album.Images[0].URL
			for _, a := range p.Item.Artists {
				artists = append(artists, a.Name)
			}
			txt = mdparser.GetUserMention(inlq.From.FirstName, inlq.From.Id).Normal(" is listening to\n").
				Italic(strings.Join(artists, ", ")).Normal(" - ").Link(p.Item.Name, p.Item.ExternalURLs["spotify"])
			title = fmt.Sprintf("%s - %s", artists, p.Item.Name)
		}
		result = gotgbot.InlineQueryResultArticle{
			Id:    uuid.NewString(),
			Title: title,
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText:           txt.ToString(),
				ParseMode:             "markdownv2",
				DisableWebPagePreview: true,
			},
			ThumbUrl: img,
		}

	}
	results = append(results, result)
	_, err = inlq.Answer(
		b, results,
		&gotgbot.AnswerInlineQueryOpts{IsPersonal: true, CacheTime: 1},
	)
	if err != nil {
		return err
	}
	return nil

}
