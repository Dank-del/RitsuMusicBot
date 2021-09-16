package handlers

const (
	startCommand          = "start"
	statusMessage         = "status"
	statusCommand         = "nowplaying"
	registerCommand       = "setusername"
	meCommand             = "me"
	topArtistsCommand     = "topartists"
	helpCommand           = "help"
	historyCommand        = "history"
	getStatusCommand      = "getstatus"
	lyricsCommand         = "lyrics"
	release               = "Beta"
	aboutCommand          = "about"
	setVisibilityCommand  = "visible"
	uploadDatabaseCommand = "uploaddb"
)

const (
	albumPrefix = "st_"
	albumText   = "Show album art"
	hideText    = "Hide album art"
)

const (
	urlRegEx = `([a-zA-Z\d]+://)?((\w+:\w+@)?([a-zA-Z\d.-]+\.[A-Za-z]{2,4})(:\d+)?(/.*)?)`
)
