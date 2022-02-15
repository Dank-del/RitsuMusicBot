package auth

import (
	"sync"
)

var (
	OauthMutex = &sync.Mutex{}
)
