package auth

// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.

import (
	"fmt"
	"github.com/zmb3/spotify/v2/auth"
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	"log"
	"net/http"
	"strconv"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.

func SpotifyAuthServer() {
	// first start an HTTP server
	log.Println("Listening on:", config.Local.Config.ServerAddr)
	config.Local.SpotifyAuthenticator = spotifyauth.New(spotifyauth.WithRedirectURL(config.Local.Config.SpotifyRedirectUri),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(config.Local.Config.SpotifyClientID),
		spotifyauth.WithClientSecret(config.Local.Config.SpotifyClientSecret))
	// log.Println("Please log in to Spotify by visiting the following page in your browser:", SpotifyAuthUrl)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	err := http.ListenAndServe(config.Local.Config.ServerAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	st := r.FormValue("state")
	tok, err := config.Local.SpotifyAuthenticator.Token(r.Context(), st, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		logging.SUGARED.Error(err)
		return
	}
	userId, err := strconv.ParseInt(st, 10, 64)
	if err != nil {
		logging.SUGARED.Error(err)
		return
	}
	database.UpdateSpotifyUser(userId, tok.RefreshToken)
	http.Redirect(w, r, fmt.Sprintf("https://t.me/%s?start=sp", config.Local.Bot.Username), 301)
	/* use the token to get an authenticated client
		htmlMarkup := fmt.Sprintf(`
	<!-- owo -->
	<body style="background-color:black;">
	<h3 style="color:#FFFFFF;" > %s Authentication for  <span style="background-color: #1cac4c; color: white; padding: 0 3px;">Spotify</span> is now complete.</h3>
	<p><strong style="color:#FFFFFF;" >Copy the below text and send to @%s.</strong></p>
	<code style="color:#FFFFFF;">/start spotifyCode%s</code>
	`, config.Local.Bot.FirstName, config.Local.Bot.Username, tok.RefreshToken)
		_, err = fmt.Fprint(w, htmlMarkup)
		if err != nil {
			log.Println(err)
			return
		}
		// log.Println(client)
		// log.Println(tok)
	*/
}

/*
func GetAsToken(user *SpotifyUser) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		TokenType:    user.TokenType,
		Expiry:       user.Expiry,
	}
}
*/
