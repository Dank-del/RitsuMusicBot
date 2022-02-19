package auth

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"golang.org/x/oauth2"
)

func GetSpotifyClient(usr *database.SpotifyUser) *spotify.Client {
	ctx := context.Background()
	userAuth := spotify.New(config.Local.SpotifyAuthenticator.Client(ctx, &oauth2.Token{RefreshToken: usr.RefreshToken}))
	return userAuth
}
