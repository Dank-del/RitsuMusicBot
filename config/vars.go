package config

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gotgbot/ratelimiter/ratelimiter"
)

var Data *DaemonConfig
var mdMessageOpt = &gotgbot.SendMessageOpts{
	ParseMode:                "markdownv2",
	AllowSendingWithoutReply: true,
}

var Limiter *ratelimiter.Limiter
