package database

import "sync"

var chatsMap = map[int64]*Chat{}
var usersMap = map[int64]*User{}
var botUserMapById = map[int64]*BotUser{}
var botUserMapByUsername = map[string]*BotUser{}
var databaseMutex = &sync.RWMutex{}
