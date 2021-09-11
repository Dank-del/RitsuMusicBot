package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var Data *DaemonConfig
var mdMessageOpt = &gotgbot.SendMessageOpts{
	ParseMode:                "markdownv2",
	AllowSendingWithoutReply: true,
}

func GetDefaultMdOpt() *gotgbot.SendMessageOpts {
	if mdMessageOpt.ReplyToMessageId != 0 {
		mdMessageOpt.ReplyToMessageId = 0
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
		fmt.Println("Error in parsing config:", err)
	}

	return err
}
