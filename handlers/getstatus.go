package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
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
			_, err = msg.Reply(b, err.Error(), nil)
			return err
		} else if d == nil {
			_, err = msg.Reply(b, "No status found for this user", nil)
			return err
		} else {
			if config.Data.IsSudo(msg.From.Id) {
				status := config.Limiter.GetStatus(user.Id)
				if status != nil && status.IsLimited() {
					d = d.AppendItalic("\n\n>This user is limited since: " +
						status.Last.String())
				}
			}
			_, err := b.SendMessage(msg.Chat.Id, d.ToString(), &gotgbot.SendMessageOpts{
				ParseMode:             "markdownv2",
				DisableWebPagePreview: true,
			})
			return err
		}
	} else {
		_, err := b.SendMessage(msg.Chat.Id, "<i>Command should be used in reply to a person.</i>",
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			})
		return err
	}
}
