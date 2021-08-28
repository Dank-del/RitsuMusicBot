package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"html"
)

func helpHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	chat := ctx.Message.Chat
	_, err := b.SendChatAction(chat.Id, "typing")
	if err != nil {
		return err
	}
	txt := fmt.Sprintf("<b>Hi, I'm %s</b>\n<i>I let you flex your last.FM on telegram\n\n</i>", html.EscapeString(b.FirstName))
	txt += "<b>Available commands</b>\n"
	txt += fmt.Sprintf("/%s - starts the bot.\n", startCommand)
	txt += fmt.Sprintf("/%s - makes me send this message.\n", helpCommand)
	txt += fmt.Sprintf("/%s - register yourself on the bot.\n", registerCommand)
	txt += fmt.Sprintf("/%s - makes me send a list of tracks you recently played.\n", historyCommand)
	txt += fmt.Sprintf(`/%s - show recently played track, also works on saying "status"`, statusCommand)
	txt += "\n\n<b>Built with ❤ by Sayan Biswas (2021)</b>"

	if chat.Type == "private" {
		_, err := msg.Reply(b, txt, &gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	} else {
		_, err = b.SendMessage(chat.Id, "<i>Command only for PM</i>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	return nil
}
