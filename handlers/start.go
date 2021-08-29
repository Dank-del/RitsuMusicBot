package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"html"
)

func startHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	chat := msg.Chat
	var message string
	if chat.Type == "private" {
		message = fmt.Sprintf("<b>Hi %s, I'm %s</b>!\n<i>Run /help to learn more!</i>",
			html.EscapeString(user.FirstName),
			html.EscapeString(b.FirstName))
	} else {
		message = fmt.Sprintf("Hi, I'm %s, ready to show off your music", b.FirstName)
	}
	_, err := msg.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "html"})
	if err != nil {
		return err
	}
	return nil
}
