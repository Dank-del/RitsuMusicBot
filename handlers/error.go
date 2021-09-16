package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HasteBin struct {
	Key string `json:"key"`
}

func postError(error string) (res *HasteBin, err error) {
	url := "https://hastebin.com/documents"
	resp, err := http.Post(url, "text/plain", strings.NewReader(error))
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return nil, fmt.Errorf("unable to post document: %v", err)
	}

	if resp.StatusCode != 200 {
		logging.SUGARED.Error(resp.StatusCode)
		return nil, errors.New("Server responded with error " + resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return nil, fmt.Errorf("unable to read reply contents: %v", err)
	}
	h := new(HasteBin)
	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil

}

var errorHandler = func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {

	chat := ctx.EffectiveChat
	tgErr := err.(*gotgbot.TelegramError)

	// these two just makes sure that errors are not logged and passed
	// as these are predefined by library
	if err == ext.ContinueGroups {
		return ext.DispatcherActionContinueGroups
	}

	if err == ext.EndGroups {
		return ext.DispatcherActionEndGroups
	}

	// if bot is not able to send any message to chat, it will leave the group
	if tgErr.Description == "Bad Request: have no rights to send a message" {
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

	// Generate a new Sha1 Hash
	shaHash := sha512.New()
	shaHash.Write([]byte(string(errorJson) + string(updateJson)))
	hash := hex.EncodeToString(shaHash.Sum(nil))

	logUrl, err := postError(string(errorJson) + "\n\n\n" + string(updateJson) + "\n\n\n" + tgErr.Error()) // helpers.CreateTelegraphPost("Error Report", string(errorJson)+"<br><br>"+string(updateJson)+"<br><br>"+tgErr.Error())
	msg := mdparser.GetBold("⚠️ An ERROR Occurred ⚠️").AppendNormal("\n\n")
	msg = msg.AppendNormal("An exception was raised while handling an update.").AppendNormal("\n\n")
	msg = msg.AppendBold("Error ID").AppendNormal(": ").AppendMono(hash).AppendNormal("\n")
	msg = msg.AppendBold("Chat ID").AppendNormal(": ").AppendMono(strconv.FormatInt(uMsg.Chat.Id, 10)).AppendNormal("\n")
	var tmpmarkup gotgbot.InlineKeyboardButton
	keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
	if err != nil {
		// logging.SUGARED.Error(err.Error())
		// msg = msg.AppendBold("Error Log").AppendNormal(": ").AppendNormal(err.Error()).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Hastebin",
			Url:  "https://hastebin.com/",
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	} else {
		// msg = msg.AppendBold("Error Log").AppendNormal(": ").AppendNormal("https://hastebin.com/" + logUrl.Key).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Hastebin",
			Url:  "https://hastebin.com/" + logUrl.Key,
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
			logging.SUGARED.Error(err.Error())
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

var panicHandler = func(b *gotgbot.Bot, ctx *ext.Context, stack []byte) {
	defer func() {
		if err := recover(); err != nil {
			logging.SUGARED.Warn("panic occurred:", err)
		}
	}()
	update := ctx.Update
	uMsg := update.Message

	// Generate a new Sha1 Hash
	shaHash := sha512.New()
	shaHash.Write([]byte(string(stack)))
	hash := hex.EncodeToString(shaHash.Sum(nil))

	// logUrl := helpers.CreateTelegraphPost("Panic Report", string(stack))
	logUrl, err := postError("Panic Report" + "\n\n\n" + string(stack))
	msg := mdparser.GetBold("⚠️ An ERROR Occurred ⚠️").AppendNormal("\n\n")
	msg = msg.AppendNormal("An exception was raised while handling an update.").AppendNormal("\n\n")
	msg = msg.AppendBold("Panic ID").AppendNormal(": ").AppendMono(hash).AppendNormal("\n")
	msg = msg.AppendBold("Chat ID").AppendNormal(": ").AppendMono(strconv.FormatInt(uMsg.Chat.Id, 10)).AppendNormal("\n")
	var tmpmarkup gotgbot.InlineKeyboardButton
	keyboard := make([][]gotgbot.InlineKeyboardButton, 1)
	if err != nil {
		// logging.SUGARED.Error(err.Error())
		// msg = msg.AppendBold("Panic Log").AppendNormal(": ").AppendNormal(err.Error()).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Hastebin",
			Url:  "https://hastebin.com/",
		}
		keyboard[0] = append(keyboard[0], tmpmarkup)
	} else {
		// msg = msg.AppendBold("Panic Log").AppendNormal(": ").AppendNormal("https://hastebin.com/" + logUrl.Key).AppendNormal("\n\n")
		tmpmarkup = gotgbot.InlineKeyboardButton{
			Text: "Hastebin",
			Url:  "https://hastebin.com/" + logUrl.Key,
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
			}})
		if err != nil {
			logging.SUGARED.Error(err.Error())
		}
	}
	/* Send Message to Log Group
	b.SendMessage(
		config.MessageDump,
		"⚠️ An ERROR Occured ⚠️\n\n"+
			"Panic Occured"+
			"\n\n"+
			fmt.Sprintf("<b>Panic ID:</b> <code>%s</code>", hash)+
			"\n"+
			fmt.Sprintf("<b>Chat ID:</b> <code>%d</code>", uMsg.Chat.Id)+
			"\n"+
			fmt.Sprintf("<b>Command:</b> <code>%s</code>", uMsg.Text)+
			"\n"+
			fmt.Sprintf("<b>Panic Log:</b> %s", logUrl)+
			"\n\n"+
			"Please Check logs ASAP!",
		parsemode.Shtml(),
	)
	log.WithFields(
		log.Fields{
			"PanicId": hash,
			"Panic":   string(stack),
			"LogURL":  logUrl,
		},
	).Error("[Dispatcher] Panic: Stack Error Occured")
	*/
}
