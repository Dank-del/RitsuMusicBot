package handlers

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/gotgbot/ratelimiter/ratelimiter"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
)

func LoadHandlers(d *ext.Dispatcher) {
	loadLimiter(d)

	startCMD := handlers.NewCommand(startCommand, startHandler)
	helpCMD := handlers.NewCommand(helpCommand, helpHandler)
	statusMsg := handlers.NewMessage(statusFilter, statusHandler)
	statusCMD := handlers.NewCommand(statusCommand, statusHandler)
	statusInl := handlers.NewInlineQuery(statusInlineFilter, statusInline)
	setUsernameCMD := handlers.NewCommand(registerCommand, setUsername)
	meCMD := handlers.NewCommand(meCommand, meHandler)
	topArtistsCMD := handlers.NewCommand(topArtistsCommand, topArtistsHandler)
	historyCMD := handlers.NewCommand(historyCommand, historyCommandHandler)
	getStatusCMD := handlers.NewCommand(getStatusCommand, getStatusHandler)
	lyricsCMD := handlers.NewCommand(lyricsCommand, lyricsHandler)
	lyricsinl := handlers.NewInlineQuery(lyricsInlineFilter, lyricsInline)
	aboutCMD := handlers.NewCommand(aboutCommand, aboutHandler)
	setVisibilityCMD := handlers.NewCommand(setVisibilityCommand, setVisibilityHandler)
	historyCB := handlers.NewCallback(historyCallBackQuery, historyCallBackResponse)
	statusCB := handlers.NewCallback(albumCallBackQuery, albumCallBackResponse)
	tdataCB := handlers.NewCallback(tDataCallBackQuery, tdataCallbackResponse)
	uploadDBcmd := handlers.NewCommand(uploadDatabaseCommand, uploadDatabase)
	linkMsg := handlers.NewMessage(msgLinkFilter, odesliMessageHandler)
	logMsg := handlers.NewMessage(logUserFilter, logUser)
	d.AddHandler(startCMD)
	d.AddHandler(helpCMD)
	d.AddHandler(statusMsg)
	d.AddHandler(statusCMD)
	d.AddHandler(setUsernameCMD)
	d.AddHandler(meCMD)
	d.AddHandler(topArtistsCMD)
	d.AddHandler(statusInl)
	d.AddHandler(historyCMD)
	d.AddHandler(getStatusCMD)
	d.AddHandler(lyricsCMD)
	d.AddHandler(lyricsinl)
	d.AddHandler(aboutCMD)
	d.AddHandler(setVisibilityCMD)
	d.AddHandler(historyCB)
	d.AddHandler(statusCB)
	d.AddHandler(tdataCB)
	d.AddHandler(uploadDBcmd)
	d.AddHandler(linkMsg)
	d.AddHandler(logMsg)
}

func loadLimiter(d *ext.Dispatcher) {
	config.Limiter = ratelimiter.NewLimiter(d, false, false)
	config.Limiter.ConsiderUser = true
	config.Limiter.IgnoreMediaGroup = true

	// 8 messages per 12 seconds
	config.Limiter.SetFloodWaitTime(8 * time.Second)
	config.Limiter.SetMaxMessageCount(12)

	if len(config.Data.SudoUsers) != 0 {
		config.Limiter.AddExceptionID(config.Data.SudoUsers...)
	}

	config.Limiter.Start()

}
