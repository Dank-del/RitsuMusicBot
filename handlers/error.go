package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
)

type MemochoRequest struct {
	Snippet string `json:"snippet"`
}

func postError(error string) (link string, err error) {
	//url, err := kv.PasteBinSession.PasteAnonymous(error, title, "", "N", "1")
	data, _ := json.Marshal(MemochoRequest{Snippet: error})
	resp, err := http.Post("https://bin.kv2.dev/", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.SUGARED.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return resp.Request.URL.String(), nil
	}
	return "", errors.New("failed to paste due to api error")
}

var ErrorHandler = func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
	chat := ctx.EffectiveChat
	logging.SUGARED.Error(err)
	tgErr, ok := err.(*gotgbot.TelegramError)

	// if bot is not able to send any message to chat, it will leave the group
	if ok && tgErr.Description == "Bad Request: have no rights to send a message" {
		_, err := b.LeaveChat(chat.Id)
		if err != nil {
			logging.SUGARED.Error(err.Error())
			return 0
		}
		return ext.DispatcherActionEndGroups
	}

	update := ctx.Update
	uMsg := update.Message
	errorJson, _ := json.MarshalIndent(tgErr, "", "  ")
	updateJson, _ := json.Marshal(update)

	// ctx.EffectiveMessage.Reply(b,
	//   "Some Error has occurred!\nIt has been reported to the developer team!",
	//   nil,
	// )

	logUrl, err := postError(string(errorJson) + "\n\n\n" + string(updateJson) + "\n\n\n" + err.Error()) // helpers.CreateTelegraphPost("Error Report", string(errorJson)+"<br><br>"+string(updateJson)+"<br><br>"+tgErr.Error())
	msg := mdparser.GetBold("⚠️ An ERROR Occurred ⚠️").AppendNormal("\n\n")
	msg = msg.AppendNormal("An exception was raised while handling an update.").AppendNormal("\n\n")
	msg = msg.AppendBold("Error ID").AppendNormal(": ").AppendMono(uuid.New().String()).AppendNormal("\n")
	msg = msg.AppendBold("Chat ID").AppendNormal(": ").AppendMono(strconv.FormatInt(uMsg.Chat.Id, 10)).AppendNormal("\n")
	var tmpmarkup gotgbot.InlineKeyboardButton
	keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
	if err != nil {
		// logging.SUGARED.Error(err.Error())
		// msg = msg.AppendBold("Error Log").AppendNormal(": ").AppendNormal(err.Error()).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Memochō",
			Url:  logUrl,
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	} else {
		// msg = msg.AppendBold("Error Log").AppendNormal(": ").AppendNormal("https://hastebin.com/" + logUrl.Key).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Memochō",
			Url:  logUrl,
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	}

	msg = msg.AppendBold("Please Check logs ASAP!")
	for _, a := range config.Data.SudoUsers {
		_, err := b.SendMessage(
			a,
			msg.ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2", ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}},
		)
		if err != nil {
			return 0
		}
	}

	// Send Message to Log Group

	/* log stuff
	log.WithFields(
		log.Fields{
			"ErrorId":       hash,
			"TelegramError": string(errorJson),
			"Update":        string(updateJson),
			"LogURL":        logUrl,
		},
	).Error(
		tgErr.Error(),
	)*/

	return ext.DispatcherActionNoop
}

var PanicHandler = func(b *gotgbot.Bot, ctx *ext.Context, i interface{}) {
	var stack []byte
	stack = debug.Stack()
	logging.SUGARED.Error(string(stack))
	defer func() {
		if err := recover(); err != nil {
			logging.SUGARED.Warn("panic occurred:", err)
		}
	}()
	update := ctx.Update
	uMsg := update.Message
	/*stack, err := GetBytes(i)
	if err != nil {
		kigLogger.Error(err)
	}*/
	ctxJson, _ := json.MarshalIndent(ctx, "", "  ")
	// logUrl := helpers.CreateTelegraphPost("Panic Report", string(stack))
	logUrl, err := postError(string(ctxJson) + "\n\n" + fmt.Sprintf("%s\n\n%v", string(stack), i))
	msg := mdparser.GetBold("⚠️ An ERROR Occurred ⚠️").AppendNormal("\n\n")
	msg = msg.AppendNormal("An exception was raised while handling an update.").AppendNormal("\n\n")
	msg = msg.AppendBold("Panic ID").AppendNormal(": ").AppendMono(uuid.New().String()).AppendNormal("\n")
	msg = msg.AppendBold("Chat ID").AppendNormal(": ").AppendMono(strconv.FormatInt(uMsg.Chat.Id, 10)).AppendNormal("\n")
	var tmpmarkup gotgbot.InlineKeyboardButton
	keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
	if err != nil {
		// logging.SUGARED.Error(err.Error())
		// msg = msg.AppendBold("Panic Log").AppendNormal(": ").AppendNormal(err.Error()).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Memochō",
			Url:  logUrl,
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	} else {
		// msg = msg.AppendBold("Panic Log").AppendNormal(": ").AppendNormal("https://hastebin.com/" + logUrl.Key).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Memochō",
			Url:  logUrl,
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	}
	msg = msg.AppendBold("Please Check logs ASAP!")
	for _, a := range config.Data.SudoUsers {
		_, err := b.SendMessage(a, msg.ToString(),
			&gotgbot.SendMessageOpts{ParseMode: "markdownv2", ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}})
		if err != nil {
			return
		}
	}
	logging.SUGARED.Debug("recover complete")
}
