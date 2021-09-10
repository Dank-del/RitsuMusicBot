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
	getStatusCMD := handlers.NewCommand(getstatusCommand, getStatusHandler)
	lyricsCMD := handlers.NewCommand(lyricsCommand, lyricsHandler)
	lyricsinl := handlers.NewInlineQuery(lyricsInlineFilter, lyricsInline)
	aboutCMD := handlers.NewCommand(aboutCommand, aboutHandler)
	historyCB := handlers.NewCallback(historyCallBackQuery, historyCallBackResponse)
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
	d.AddHandler(historyCB)
}
