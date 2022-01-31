package handlers

import (
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/gotgbot/ratelimiter/ratelimiter"
)

func LoadHandlers(d *ext.Dispatcher) {
	loadLimiter(d)
	startCMD := handlers.NewCommand(startCommand, startHandler)
	helpCMD := handlers.NewCommand(helpCommand, helpHandler)
	statusMsg := handlers.NewMessage(statusFilter, statusHandler)
	statusCMD := handlers.NewCommand(statusCommand, statusHandler)
	changeStatusCMD := handlers.NewCommand(changeStatusCommand, changeStatusHandler)
	statusInl := handlers.NewInlineQuery(statusInlineFilter, statusInline)
	setUsernameCMD := handlers.NewCommand(registerCommand, setUsername)
	meCMD := handlers.NewCommand(meCommand, meHandler)
	gitpullCMD := handlers.NewCommand("gitpull", gitPull)
	topArtistsCMD := handlers.NewCommand(topArtistsCommand, topArtistsHandler)
	historyCMD := handlers.NewCommand(historyCommand, historyCommandHandler)
	getStatusCMD := handlers.NewCommand(getStatusCommand, getStatusHandler)
	lyricsCMD := handlers.NewCommand(lyricsCommand, lyricsHandler)
	lyricsinl := handlers.NewInlineQuery(lyricsInlineFilter, lyricsInline)
	aboutCMD := handlers.NewCommand(aboutCommand, aboutHandler)
	setVisibilityCMD := handlers.NewCommand(setVisibilityCommand, setVisibilityHandler)
	historyCB := handlers.NewCallback(historyCallBackQuery, historyCallBackResponse)
	statusCB := handlers.NewCallback(albumCallBackQuery, albumCallBackResponse)
	uploadDBcmd := handlers.NewCommand(uploadDatabaseCommand, uploadDatabase)
	linkMsg := handlers.NewMessage(msgLinkFilter, odesliMessageHandler)
	logMsg := handlers.NewMessage(logUserFilter, logUser)
	d.AddHandler(handlers.NewCommand(linkDetectCommand, setLinkDetection))
	d.AddHandler(startCMD)
	d.AddHandler(helpCMD)
	d.AddHandler(gitpullCMD)
	d.AddHandler(statusMsg)
	d.AddHandler(statusCMD)
	d.AddHandler(changeStatusCMD)
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
	d.AddHandler(uploadDBcmd)
	d.AddHandler(linkMsg)
	d.AddHandler(logMsg)
}

func loadLimiter(d *ext.Dispatcher) {
	config.Limiter = ratelimiter.NewLimiter(d, false, false)
	config.Limiter.ConsiderUser = true
	config.Limiter.IgnoreMediaGroup = true

	// 14 messages per 6 seconds
	config.Limiter.SetFloodWaitTime(6 * time.Second)
	config.Limiter.SetMaxMessageCount(14)

	if len(config.Data.SudoUsers) != 0 {
		config.Limiter.AddExceptionID(config.Data.SudoUsers...)
	}

	config.Limiter.Start()

}
