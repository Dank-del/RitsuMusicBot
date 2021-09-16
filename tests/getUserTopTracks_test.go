package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"testing"
)

func TestTopTracksOfUser(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
		return
	}
	d, err := last_fm.GetTopTracks("airi_sakura")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(d.Toptracks.Track)
}
