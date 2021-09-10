package tests

import (
	"testing"

	"gitlab.com/Dank-del/lastfm-tgbot/config"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
)

func TestGetRecents(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
	}
	d, err := last_fm.GetRecentTracksByUsername("airi_sakura", 10)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(d.Recenttracks)
}
