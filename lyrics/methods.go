package lyrics

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"golang.org/x/net/html"
)

func search(query string) (res *SearchResult, err error) {
	reqUrl := baseUrl + searchMeth + url.QueryEscape(query)
	// fmt.Println(reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Data.GeniusApiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//log.Println(string(b))

	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// parseLyrics returns the lyrics from a Genius page, or an empty []string if it couldn't find them.
func parseLyrics(r io.Reader) (lyrics []string) {
	z := html.NewTokenizer(r)

	var inLyrics bool
	var start bool
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		name, hasAttr := z.TagName()
		if string(name) == "div" && hasAttr {
			k, v, _ := z.TagAttr()
			if string(k) == "class" && string(v) == "lyrics" {
				inLyrics = true
				continue
			}
		}

		if inLyrics {
			text := string(z.Text())
			if text == "sse" {
				start = true
				continue
			}
			if text == "/sse" {
				break
			}
			if start {
				text = strings.TrimSpace(text)
				if len(text) > 0 {
					lyrics = append(lyrics, text)
				}
			}
		}
	}
	return
}

func lyricsFromUrl(url string) ([]string, error) {
	var l []string
	r, err := http.Get(url)
	if err != nil {
		return l, err
	} else if r.StatusCode == 404 {
		return l, fmt.Errorf("404: %s not found", url)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(r.Body)

	l = parseLyrics(r.Body)
	if len(l) == 0 {
		err = fmt.Errorf("no lyrics parsed")
	}
	return l, err
}

func GetLyrics(query string) ([]string, error) {
	d, err := search(query)
	if err != nil {
		return nil, err
	}

	// fmt.Println(d.Response.Hits[0].Result.URL)
	l, err := lyricsFromUrl(d.Response.Hits[0].Result.URL)
	if err != nil {
		return nil, err
	}
	return l, err

}
