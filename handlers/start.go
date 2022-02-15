package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/auth"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"golang.org/x/net/context"
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
		if len(args) > 1 && strings.HasPrefix(args[1], "sp") {
			auth.OauthMutex.Lock()
			defer auth.OauthMutex.Unlock()
			sptUsr, err := database.GetSpotifyUser(user.Id)
			if err != nil {
				message = err.Error()
			} else {
				currentUser, err := sptUsr.CurrentUser(ctxBg)
				if err != nil {
					message = err.Error()
				} else {
					message = fmt.Sprintf("Logged in as %s (%s)", currentUser.DisplayName, currentUser.Product)
				}
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
