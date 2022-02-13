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
	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2/auth"
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"log"
	"net/http"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.

func SpotifyAuthServer() {
	// first start an HTTP server
	defer log.Println("Listening on:", config.Local.Config.ServerAddr)
	SpotifyAuthenticator = spotifyauth.New(spotifyauth.WithRedirectURL(config.Local.Config.SpotifyRedirectUri),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(config.Local.Config.SpotifyClientID),
		spotifyauth.WithClientSecret(config.Local.Config.SpotifyClientSecret))
	SpotifyAuthUrl = SpotifyAuthenticator.AuthURL(state)
	log.Println("Please log in to Spotify by visiting the following page in your browser:", SpotifyAuthUrl)
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
	tok, err := SpotifyAuthenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	key := uuid.NewString()
	TokenMutex.RLock()
	TokenMap[key] = tok.RefreshToken
	TokenMutex.RUnlock()
	http.Redirect(w, r, fmt.Sprintf("https://t.me/%s?start=spotifyCode%s", config.Local.Bot.Username, key), 301)
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
