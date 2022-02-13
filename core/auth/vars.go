package auth

import (
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"sync"
)

var (
	SpotifyAuthenticator = &spotifyauth.Authenticator{}
	state                = "abc123"
	SpotifyAuthUrl       = ""
	TokenMap             = map[string]string{}
	TokenMutex           = &sync.RWMutex{}
)
