package tests

import (
	"log"
	"testing"

	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/lyrics"
)

// TestLyrics function will test the `lyrics.GetLyrics` function
// and will get lyrics of `Loote - she's all yours` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Error(err.Error())
		return
	}

	l, err := lyrics.GetLyrics("Loote - she's all yours")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(l) < 5 {
		t.Error("length of the string array is too short")
		return
	}
	t.Log(l)
}

// TestLyrics2 function will test the `lyrics.GetLyrics` function
// and will get lyrics of `LiSA - unlasting` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics2(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Errorf("got an error when tried to get config: %v", err)
		return
	}

	strs, err := lyrics.GetLyrics("LiSA - unlasting")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	if len(strs) < 5 {
		t.Error("length of the string array is too short")
		return
	}

	log.Println(strs)
}

// TestLyrics function will test the `lyrics.GetLyrics` function
// and will get lyrics of `Loote - she's all yours` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics3(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Errorf("got an error when tried to get config: %v", err)
		return
	}

	// "Machine Gun Kelly - I Think I'm OKAY
	strs, err := lyrics.GetLyrics("Machine Gun Kelly - I Think I'm OKAY (with YUNGBLUD & Travis Barker)")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	if len(strs) < 5 {
		t.Error("length of the string array is too short")
		return
	}

	log.Println(strs)
}

// TestLyrics function will test the `lyrics.GetLyrics` function
// and will get lyrics of `Loote - she's all yours` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics4(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Errorf("got an error when tried to get config: %v", err)
		return
	}

	//Feverkin - Overthought
	// "Machine Gun Kelly - I Think I'm OKAY
	strs, err := lyrics.GetLyrics("LiSA - unlasting")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	if len(strs) < 5 {
		t.Error("length of the string array is too short")
		return
	}

	log.Println(strs)
}

// TestLyrics function will test the `lyrics.GetLyrics` function
// and will get lyrics of `Loote - she's all yours` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics5(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Errorf("got an error when tried to get config: %v", err)
		return
	}

	// "Machine Gun Kelly - I Think I'm OKAY
	strs, err := lyrics.GetLyrics("Feverkin - Overthought")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	if len(strs) < 5 {
		t.Error("length of the string array is too short")
		return
	}

	log.Println(strs)
}

// TestLyrics function will test the `lyrics.GetLyrics` function
// and will get lyrics of `Loote - she's all yours` song.
// final Test results: PASS
// last date: Sep 9 2021;
func TestLyrics6(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		t.Errorf("got an error when tried to get config: %v", err)
		return
	}

	// "Machine Gun Kelly - I Think I'm OKAY
	strs, err := lyrics.GetLyrics("bo en - my time")
	if err != nil {
		t.Errorf("got an error from GetLyrics func: %v", err)
		return
	}

	if len(strs) < 5 {
		t.Error("length of the string array is too short")
		return
	}

	log.Println(strs)
}
