package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher) {
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
	statusCB := handlers.NewCallback(statusCallBackQuery, statusCallBackResponse)
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
	d.AddHandler(logMsg)
}
