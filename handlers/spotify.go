package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func linkSpotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	if chat.Type != "private" {
		_, err := ctx.EffectiveMessage.Reply(b, "Try that in pm", nil)
		return err
	}

	return ext.EndGroups
}
