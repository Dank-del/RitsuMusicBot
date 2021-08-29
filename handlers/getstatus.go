package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func getStatusHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	if msg.ReplyToMessage != nil {
		user := msg.ReplyToMessage.From
		d, err := getStatus(user)
		if err != nil {
			_, err := b.SendMessage(msg.Chat.Id, err.Error(), &gotgbot.SendMessageOpts{})
			return err
		} else {
			_, err := b.SendMessage(msg.Chat.Id, d, &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
			return err
		}
	} else {
		_, err := b.SendMessage(msg.Chat.Id, "<i>Command to be used in reply to a person.</i>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}
}
