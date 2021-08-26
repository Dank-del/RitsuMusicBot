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
	setUsernameCMD := handlers.NewCommand(registerCommand, setUsername)
	meCMD := handlers.NewCommand(meCommand, meHandler)
	topArtistsCMD := handlers.NewCommand(topArtistsCommand, topArtistsHandler)
	d.AddHandler(startCMD)
	d.AddHandler(helpCMD)
	d.AddHandler(statusMsg)
	d.AddHandler(statusCMD)
	d.AddHandler(setUsernameCMD)
	d.AddHandler(meCMD)
	d.AddHandler(topArtistsCMD)

}
