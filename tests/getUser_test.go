package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
	"testing"
)

func TestGetUser(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
	}
	d, err := last_fm.GetLastFMUser("airi_sakura")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(d.User.Image[0].Text)
}
