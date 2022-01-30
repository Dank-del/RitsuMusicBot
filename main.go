package main

import (
	"fmt"
	"github.com/Dank-del/MusixScrape/musixScrape"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"gitlab.com/Dank-del/lastfm-tgbot/handlers"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	loggerMgr := logging.InitZapLog()
	zap.ReplaceGlobals(loggerMgr)
	defer func(loggerMgr *zap.Logger) {
		err := loggerMgr.Sync()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(loggerMgr) // flushes buffer, if any
	Logger := loggerMgr.Sugar()
	logging.SUGARED = loggerMgr.Sugar()
	err := config.GetConfig()
	if err != nil {
		Logger.Error(err.Error())
	}
	undo := zap.RedirectStdLog(loggerMgr)
	defer undo()
	config.Local.Config = config.Data
	config.Local.MusixMatchSession = musixScrape.New(nil)
	logging.SUGARED.Info("Starting daemon..")
	b, err := gotgbot.NewBot(config.Data.BotToken, &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	logging.SUGARED.Info(fmt.Sprintf("GetTimeout %s", gotgbot.DefaultGetTimeout))
	logging.SUGARED.Info(fmt.Sprintf("PostTimeout %s", gotgbot.DefaultPostTimeout))
	if err != nil {
		logging.SUGARED.Error(err.Error())
	}
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		DispatcherOpts: ext.DispatcherOpts{
			MaxRoutines: -1,
			Error:       handlers.ErrorHandler,
			Panic:       handlers.PanicHandler,
		},
		ErrorLog: zap.NewStdLog(loggerMgr),
	})
	dispatcher := updater.Dispatcher
	handlers.LoadHandlers(dispatcher)
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		logging.SUGARED.Error(fmt.Sprintf("Failed to start polling due to %s", err.Error()))
	}
	database.StartDatabase(b.Id)
	logging.SUGARED.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))
	updater.Idle()
}
