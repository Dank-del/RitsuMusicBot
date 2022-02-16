package handlers

import (
	"fmt"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	lastfm "gitlab.com/Dank-del/lastfm-tgbot/libs/last.fm"
)

func setUsername(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	args := ctx.Args()
	chatID := msg.Chat.Id
	d, err := database.GetChat(chatID)
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return ext.EndGroups
	}
	// print(d.GetStatusMessage())
	var m mdparser.WMarkDown
	if len(args) == 1 {
		m = mdparser.GetItalic("Usage: ").AppendItalic(args[0])
		m = m.AppendItalic(" airi_sakura")
		_, err := msg.Reply(b, m.ToString(), config.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	} else {
		username := args[1]
		user, err := lastfm.GetLastFMUser(username)
		if err != nil || user == nil || user.User == nil {
			m = mdparser.GetItalic("\n> check the spelling and try again")
			_, err := msg.Reply(b, m.ToString(), config.GetDefaultMdOpt())
			if err != nil {
				return err
			}
			return nil
		}
		m = mdparser.GetItalic(`Username set as `).AppendMono(username)
		m = m.AppendItalic(fmt.Sprintf(`, enjoy flexing "%s"`, d.GetStatusMessage()))
		database.UpdateLastFMUserInDB(msg.From.Id, username)
		database.UpdateBotUser(msg.From.Id, msg.From.Username, false)
		_, err = msg.Reply(b, m.ToString(), config.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	}
	return nil
}

func linkSpotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "private" {
		_, err := ctx.EffectiveMessage.Reply(b,
			mdparser.GetBold("Send me ").Mono("/"+linkSpotifyCommand).Bold(" to authenticate").ToString(),
			config.GetDefaultMdOpt())
		return err
	}
	config.LinkMutex.RLock()
	defer config.LinkMutex.RUnlock()
	uid := uuid.NewString()
	config.LinkMap[uid] = ctx.EffectiveUser.Id
	txt := mdparser.GetBold("Click ").Link("this",
		config.Local.SpotifyAuthenticator.AuthURL(uid)).
		Bold(" link to authenticate.")
	_, err := ctx.EffectiveMessage.Reply(b, txt.ToString(), config.GetDefaultMdOpt())
	return err
}

func unlinkSpotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "private" {
		_, err := ctx.EffectiveMessage.Reply(b, "Try that in my pm.", nil)
		return err
	}
	var msg string
	_, err := database.GetSpotifyUser(ctx.EffectiveUser.Id)
	if err != nil {
		msg = err.Error()
	} else {
		database.DelSpotifyUser(ctx.EffectiveUser.Id)
		msg = fmt.Sprintf("Unlink successful.")
	}
	_, err = ctx.EffectiveMessage.Reply(b, msg, nil)
	return err
}
