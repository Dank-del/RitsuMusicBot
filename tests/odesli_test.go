package tests

import (
	"gitlab.com/Dank-del/lastfm-tgbot/libs/odesli"
	"testing"
)

func TestOdesliByUrl(t *testing.T) {
	r, err := odesli.GetLinks("https://open.spotify.com/track/7iYMW4Ww40RAbfeccmzMQa?si=069155796edc4a61")
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(r.EntitiesByUniqueID[r.EntityUniqueID].Title)
	t.Log(r.LinksByPlatform.Tidal)

}
