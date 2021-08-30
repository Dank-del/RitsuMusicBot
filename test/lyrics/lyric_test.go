package lyric_test

import (
	"log"
	"testing"

	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/lyrics"
)

func TestLyric(t *testing.T) {
	config.GetConfig()
	strs, err := lyrics.GetLyrics("LiSA - unlasting")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	log.Println(strs)
}
