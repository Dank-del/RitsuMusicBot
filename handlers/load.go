package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)



func LoadHandlers(d *ext.Dispatcher) {
	statusCMD := handlers.NewMessage(statusFilter, statusHandler)
	setUsernameCMD := handlers.NewCommand("setusername", setUsername)
	meCMD := handlers.NewCommand("me", meHandler)
	d.AddHandler(statusCMD)
	d.AddHandler(setUsernameCMD)
	d.AddHandler(meCMD)

}
