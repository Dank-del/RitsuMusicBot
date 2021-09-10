package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"testing"
)

func TestGetLastFmTrack(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
	}
	e, err := last_fm.GetLastfmTrack("Pelican Fanclub", "ディザイア")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(e.Track)
}
