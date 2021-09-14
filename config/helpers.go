package config

import (
	"encoding/json"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"os"

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
	var file *os.File
	var err error
	if os.PathSeparator == '/' {
		// linux config
		file, err = os.Open("config.json")
	} else {
		// winhoes (I mean windows) config
		file, err = os.Open("E:\\gits\\LastFM-TG\\config.json")
	}

	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Data)
	if err != nil {
		logging.SUGARED.Errorf("Error in parsing config: %s", err)
	}

	return err
}
