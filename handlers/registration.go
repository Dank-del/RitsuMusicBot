package handlers

import (
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
)

func setUsername(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	args := ctx.Args()
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
		user, err := last_fm.GetLastFMUser(username)
		if err != nil || user == nil || user.User == nil {
			m = mdparser.GetItalic("\n> check the spelling and try again")
			_, err := msg.Reply(b, m.ToString(), config.GetDefaultMdOpt())
			if err != nil {
				return err
			}
			return nil
		}
		m = mdparser.GetItalic(`Username set as `).AppendMono(username)
		m = m.AppendItalic(`, enjoy flexing "status"`)
		database.UpdateLastFMUserInDB(msg.From.Id, username)
		database.UpdateBotUser(msg.From.Id, msg.From.Username, false)
		_, err = msg.Reply(b, m.ToString(), config.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	}
	return nil
}
