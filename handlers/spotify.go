package handlers

import (
	"context"
	"fmt"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/Dank-del/MusixScrape/musixScrape"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/zmb3/spotify/v2"
	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/utilities"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"strings"
)

func spotifyNow(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	ctxBg := context.Background()
	spotifyUser, err := database.GetSpotifyUser(user.Id)
	if err != nil {
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
			_, err = b.SendPhoto(ctx.EffectiveChat.Id, img, &gotgbot.SendPhotoOpts{
				ParseMode:   "markdownv2",
				Caption:     txt.ToString(),
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard},
			})
			return err
		}
	}
	return ext.EndGroups
}
