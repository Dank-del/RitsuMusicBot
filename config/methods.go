package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Data *DaemonConfig

func GetConfig() error {
	file, err := os.Open("E:\\gits\\LastFM-TG\\config.json")
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
