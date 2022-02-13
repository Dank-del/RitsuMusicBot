package handlers

import (
	"fmt"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/zmb3/spotify/v2"
	auth2 "gitlab.com/Dank-del/lastfm-tgbot/core/auth"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"html"
	"strings"
)

func startHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	chat := msg.Chat
	args := ctx.Args()
	var message string
	if chat.Type == "private" {
		ctxBg := context.Background()
		if len(args) > 1 && strings.HasPrefix(args[1], "spotifyCode") {
			auth2.TokenMutex.RLock()
			defer auth2.TokenMutex.RUnlock()
			code := auth2.TokenMap[strings.ReplaceAll(args[1], "spotifyCode", "")]
			if code == strongStringGo.EMPTY {
				return ext.EndGroups
			}
			auth := spotify.New(auth2.SpotifyAuthenticator.Client(ctxBg, &oauth2.Token{RefreshToken: code}))
			currentUser, err := auth.CurrentUser(ctxBg)
			if err != nil {
				message = err.Error()
			} else {
				database.UpdateSpotifyUser(user.Id, code)
				message = fmt.Sprintf("Logged in as %s (%s)", currentUser.DisplayName, currentUser.Product)
			}
		} else {
			message = fmt.Sprintf("<b>Hi %s, I'm %s</b>!\n<i>Run /help to learn more!</i>",
				html.EscapeString(user.FirstName),
				html.EscapeString(b.FirstName))
		}
	} else {
		message = fmt.Sprintf("Hi, I'm %s, ready to show off your music", b.FirstName)
	}
	_, err := msg.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "html"})
	if err != nil {
		return err
	}
	return nil
}
