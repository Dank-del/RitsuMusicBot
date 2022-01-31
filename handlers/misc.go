package handlers

import (
	"fmt"
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/core/utilities"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
)

func aboutHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}

	txt := mdparser.GetBold(fmt.Sprintf("%s - %s", b.FirstName, release)).AppendNormal("\n\n")
	txt = txt.AppendItalic("Exists for the sole reason of flexing your taste of music").AppendNormal("\n\n")
	h, err := os.Hostname()
	if err != nil {
		txt = txt.AppendNormal("Node: ").AppendMono(err.Error()).AppendNormal("\n")
	} else {
		txt = txt.AppendNormal("Node: ").AppendMono(h).AppendNormal("\n")
	}

	txt = txt.AppendMono(strconv.FormatInt(database.GetLastmUserCount(), 10)).AppendNormal(" Last.FM accounts registered").AppendNormal("\n")
	txt = txt.AppendMono(strconv.FormatInt(database.GetBotUserCount(),
		10)).AppendNormal(" telegram accounts noticed by ").AppendMention(b.FirstName, b.Id).AppendNormal("\n")

	txt.Bold("Goroutines running: ").Mono(fmt.Sprintf("%v", runtime.NumGoroutine())).Normal("\n")
	txt.Bold("CPUs available: ").Mono(fmt.Sprintf("%v", runtime.NumCPU())).Normal("\n")
	txt = txt.AppendBold("Runtime: ").AppendMono(runtime.Version()).AppendNormal("\n\n")
	txt = txt.AppendBold("Built with ❤ by Sayan Biswas (2021)")

	_, err = msg.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	if err != nil {
		return err
	}
	return nil
}

func logUserFilter(msg *gotgbot.Message) bool {
	return len(msg.Text) > 0
}

func logUser(_ *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.Message.From
	dbuser, err := database.GetBotUserByID(user.Id)
	if err != nil {
		return err
	}
	if dbuser == nil {
		database.UpdateBotUser(user.Id, user.Username, false)
	} else {
		database.UpdateBotUser(user.Id, user.Username, dbuser.ShowProfile)
	}
	return ext.ContinueGroups
}

func setVisibilityHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	args := ctx.Args()
	if len(args) == 1 {
		_, err := msg.Reply(b, mdparser.GetItalic(fmt.Sprintf("Usage: %s yes/no", args[0])).ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return nil
	}

	switch v := strings.ToLower(args[1]); v {
	case "yes":
		database.UpdateBotUser(user.Id, user.Username, true)
		_, err := msg.Reply(b, mdparser.GetItalic(`Success, your profile will now be visible on "status".`).ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	case "no":
		database.UpdateBotUser(user.Id, user.Username, false)
		_, err := msg.Reply(b, mdparser.GetItalic(`Success, your profile won't be visible on "status" anymore.`).ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	default:
		_, err := msg.Reply(b, mdparser.GetItalic(`Expected Yes/No, received `).AppendMono(args[1]).ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadDatabase(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	if config2.Data.IsSudo(user.Id) {
		txt := mdparser.GetBold("Database backup for ").AppendMention(b.FirstName, b.Id).AppendNormal("\n")
		// txt = txt.AppendItalic("Exported on: ").AppendItalic(time.Now().Local().String())
		fileName := fmt.Sprintf("%d.db", b.Id)
		f, err := os.Open(fileName)
		if err != nil {
			logging.SUGARED.Error(err.Error())
			return ext.EndGroups
		}
		namedFile := gotgbot.NamedFile{
			File:     f,
			FileName: fileName,
		}

		private := msg.Chat.Type != "private" && shouldSendPrivate(strings.ToLower(msg.Text))
		if private {
			_, err = b.SendDocument(msg.From.Id, namedFile,
				&gotgbot.SendDocumentOpts{
					Caption:   txt.ToString(),
					ParseMode: "markdownv2",
				})
			if err != nil {
				warnError(msg, b, err)
				return ext.EndGroups
			}

			_, err = msg.Reply(b, mdparser.GetItalic("db backup has been sent to you.").ToString(),
				config2.GetDefaultMdOpt())
			if err != nil {
				logging.SUGARED.Errorf(err.Error())
				return ext.EndGroups
			}

		} else {
			_, err = b.SendDocument(msg.Chat.Id, namedFile,
				&gotgbot.SendDocumentOpts{
					Caption:                  txt.ToString(),
					ParseMode:                "markdownv2",
					ReplyToMessageId:         msg.MessageId,
					AllowSendingWithoutReply: false,
				})

			if err != nil {
				warnError(msg, b, err)
				return ext.EndGroups
			}
		}
	} else {
		_, err := msg.Reply(b, mdparser.GetItalic("Get the fuck away from me..").ToString(),
			config2.GetDefaultMdOpt())

		if err != nil {
			logging.SUGARED.Error(err.Error())
			return err
		}
	}
	return nil
}

func shouldSendPrivate(text string) bool {
	return strings.Contains(text, "pv") || strings.Contains(text, "pm") ||
		strings.Contains(text, "private")
}

func warnError(msg *gotgbot.Message, b *gotgbot.Bot, err error) {
	_, err = msg.Reply(b, mdparser.GetItalic(err.Error()).ToString(),
		config2.GetDefaultMdOpt())
	if err != nil {
		logging.SUGARED.Errorf(err.Error())
	}
}

func changeStatusHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	chat := msg.Chat
	args := ctx.Args()

	if !utilities.IsUserAdmin(b, &chat, user.Id) {
		_, err := msg.Reply(b, mdparser.GetItalic("You are not authorized to do that!").ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			logging.SUGARED.Error(err.Error())
			return err
		}
		return nil
	}

	if len(args) == 1 {
		_, err := msg.Reply(b, mdparser.GetItalic(fmt.Sprintf("Usage: %s <status>", args[0])).ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return nil
	}

	if len(args) > 1 {
		status := strings.Join(args[1:], " ")
		database.UpdateChat(chat.Id, status)
		_, err := msg.Reply(b, mdparser.GetBold(`Success, `).AppendNormal("status message was updated").ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			return err
		}
	}
	return nil
}

func setLinkDetection(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	chat := msg.Chat
	args := ctx.Args()
	data, err := database.GetChat(chat.Id)
	if err != nil {
		return err
	}
	if !utilities.IsUserAdmin(b, &chat, user.Id) {
		_, err := msg.Reply(b, mdparser.GetItalic("You are not authorized to do that!").ToString(),
			config2.GetDefaultMdOpt())
		if err != nil {
			logging.SUGARED.Error(err.Error())
			return err
		}
		return nil
	}

	if len(args) == 1 {
		_, err := msg.Reply(b, mdparser.GetItalic(fmt.Sprintf("Usage: %s <on/off/yes/no>", args[0])).ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return nil
	}
	enable := strongStringGo.ToBool(args[1])
	data.SetLinkDetection(enable)
	if !enable {
		_, err := msg.Reply(b, mdparser.GetItalic(fmt.Sprintf("Disabled link detection in %s", chat.Title)).ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return ext.EndGroups
	} else {
		_, err := msg.Reply(b, mdparser.GetItalic(fmt.Sprintf("Enabled link detection in %s", chat.Title)).ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return ext.EndGroups
	}
}

func gitPull(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	usr := ctx.EffectiveUser
	if !config2.Data.IsSudo(usr.Id) {
		_, err := msg.Reply(b, mdparser.GetItalic("You are not a sudo user.").ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
		if err != nil {
			return err
		}
		return nil
	}
	a := exec.Command("git", "reset", "--hard")
	_, err := a.Output()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull")
	m, err := cmd.Output()
	if err != nil {
		return err
	}
	_, err = msg.Reply(b, mdparser.GetMono(string(m)).ToString(), &gotgbot.SendMessageOpts{ParseMode: "markdownv2"})
	if err != nil {
		return err
	}

	exec.Command("cmd.exe", "run.bat")
	os.Exit(1)
	return nil
}
