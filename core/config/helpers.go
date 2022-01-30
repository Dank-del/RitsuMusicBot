package config

import (
	"encoding/json"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"io/ioutil"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// GetDefaultMdOpt function returns the default `*gotgbot.SendMessageOpts`
// value of the application.
// it won't allocate a new struct, so it will reduce the memory usage
// by a great amount.
func GetDefaultMdOpt() *gotgbot.SendMessageOpts {
	if mdMessageOpt.ReplyToMessageId != 0 {
		mdMessageOpt.ReplyToMessageId = 0
	}

	if mdMessageOpt.Entities != nil {
		mdMessageOpt.Entities = nil
	}

	if mdMessageOpt.ReplyMarkup != nil {
		mdMessageOpt.ReplyMarkup = nil
	}

	return mdMessageOpt
}

func GetConfig() error {
	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Data)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		logging.SUGARED.Errorf("Error in parsing config: %s", err)
	}

	return err
}
