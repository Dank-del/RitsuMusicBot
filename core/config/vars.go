package config

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gotgbot/ratelimiter/ratelimiter"
	"sync"
)

var Data *DaemonConfig
var mdMessageOpt = &gotgbot.SendMessageOpts{
	ParseMode:                "markdownv2",
	AllowSendingWithoutReply: true,
}

var Limiter *ratelimiter.Limiter
var Local = &DaemonLocal{}
var LinkMap = map[string]int64{}
var LinkMutex = &sync.RWMutex{}
