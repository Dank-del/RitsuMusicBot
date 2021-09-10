package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"gitlab.com/Dank-del/lastfm-tgbot/handlers"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"net/http"
)

func main() {
	err := config.GetConfig()
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Info("Starting daemon..")
	b, err := gotgbot.NewBot(config.Data.BotToken, &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	logging.Info(fmt.Sprintf("GetTimeout %s", gotgbot.DefaultGetTimeout))
	logging.Info(fmt.Sprintf("PostTimeout %s", gotgbot.DefaultPostTimeout))
	if err != nil {
		logging.Error(err.Error())
	}
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher
	handlers.LoadHandlers(dispatcher)
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to start polling due to %s", err.Error()))
	}
	database.StartDatabase(b.Id)
	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))
	updater.Idle()
}
