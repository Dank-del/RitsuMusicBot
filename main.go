package main

import (
	"fmt"
	"github.com/Dank-del/MusixScrape/musixScrape"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/auth"
	config2 "gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"gitlab.com/Dank-del/lastfm-tgbot/handlers"
	"gitlab.com/toby3d/telegraph"
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
	err := config2.GetConfig()
	if err != nil {
		Logger.Error(err.Error())
	}
	undo := zap.RedirectStdLog(loggerMgr)
	defer undo()
	config2.Local.Config = config2.Data
	config2.Local.MusixMatchSession = musixScrape.New()
	logging.SUGARED.Info("Starting daemon..")
	b, err := gotgbot.NewBot(config2.Data.BotToken, &gotgbot.BotOpts{
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
	tf, err := telegraph.CreateAccount(telegraph.Account{
		AuthorName: b.User.Username,
		ShortName:  b.User.FirstName,
		AuthorURL:  fmt.Sprintf("https://telegram.dog/%s", b.User.Username),
	})
	if err != nil {
		log.Fatalf("Error while registering on telegra.ph: %s", err.Error())
		return
	}
	config2.Local.Bot = b
	config2.Local.TelegraphClient = tf
	database.StartDatabase(b.Id)
	auth.SpotifyAuthServer()
	logging.SUGARED.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))
	updater.Idle()
}
