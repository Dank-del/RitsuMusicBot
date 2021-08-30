package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/lyrics"
	"testing"
)

func TestLyrics(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
	}
	l, err := lyrics.GetLyrics("Loote - she's all yours")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(l)
}
