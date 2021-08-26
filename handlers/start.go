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

	message := fmt.Sprintf("<b>Hi %s, I'm %s</b>!\n<i>Run /help to learn more!</i>", html.EscapeString(user.FirstName), html.EscapeString(b.FirstName))

	_, err := msg.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "html"})
	if err != nil {
		return err
	}
	return nil
}
