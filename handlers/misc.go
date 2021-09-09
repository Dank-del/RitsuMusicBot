package handlers

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func aboutHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}

	txt := mdparser.GetBold(fmt.Sprintf("%s - %s", b.FirstName, release)).AppendNormal("\n\n")
	txt = txt.AppendItalic("Exists for the sole reason of flexing your taste of music").AppendNormal("\n")
	h, err := os.Hostname()
	if err != nil {
		txt = txt.AppendNormal("Node: ").AppendMono(err.Error()).AppendNormal("\n")
	} else {
		txt = txt.AppendNormal("Node: ").AppendMono(h).AppendNormal("\n")
	}
	txt = txt.AppendBold("Runtime: ").AppendMono(runtime.Version()).AppendNormal("\n\n")
	txt = txt.AppendBold("Built with ❤ by Sayan Biswas (2021)")

	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	if err != nil {
		return err
	}
	return nil
}
