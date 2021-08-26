package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"testing"
)

func TestGetTopArtists(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
	}
	d, err := last_fm.GetTopArtists()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(d.Artists.Artist)
}
