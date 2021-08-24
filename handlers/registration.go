package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
)

func setUsername(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	args := ctx.Args()
	if len(args) == 1 {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Usage: %s airi_sakura</i>", args[0]), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	} else {
		username := args[1]
		database.UpdateUser(msg.From.Id, username)
		_, err := msg.Reply(b, fmt.Sprintf(`<i>Username set as %s, enjoy flexing "status"</i>`, username), &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return err
		}
	}
	return nil
}