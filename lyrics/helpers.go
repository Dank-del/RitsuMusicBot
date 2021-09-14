package lyrics

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"unicode"

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
			logging.SUGARED.Error(err.Error())
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

// parseLyrics returns the lyrics from a Genius page, or an empty []string
// if it couldn't find them.
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
		return nil, err
	} else if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		} else if r.StatusCode == http.StatusInternalServerError {
			return nil, ErrAPIDownError
		}

		return nil, ErrAPINotAvailable
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.SUGARED.Error(err.Error())
		}
	}(r.Body)

	l = parseLyrics(r.Body)
	if len(l) == 0 {
		return nil, ErrNotFound
	}

	return l, nil
}

func GetLyrics(query string) ([]string, error) {
	d, err := search(query)
	if err != nil {
		return nil, err
	}

	best := getHit(query, d.Response.Hits)
	if best == nil {
		return nil, ErrGoodResultNotFound // TODO: err
	}

	// fmt.Println(d.Response.Hits[0].Result.URL)
	l, err := lyricsFromUrl(best.URL)
	if err != nil {
		return nil, err
	}

	return l, err
}

// smartFind function will super cool algorithms will check
// and see if the user query string matches with the search
// title or not.
// please notice that you don't need to use `strings.ToLower`
// function for `userQuery` argument, since this value should
// be lowered by `getHit` function for reducing the running time.
func smartFind(userQuery, searchTitle string, en bool) int {
	var point int

	searchTitle = strings.ToLower(searchTitle)
	if strings.Contains(searchTitle, instrumentalValue) {
		// Instrumental version of the song is meh...
		// in the first place, how is it instrumental version
		// and it has lyric??
		point -= 2
	}

	queries := strings.Split(userQuery, " ")
	for _, q := range queries {
		if len(q) < 2 && q != "i" {
			continue
		} else if q == unimportant1 || q == unimportant2 ||
			q == unimportant3 || q == unimportant4 {
			continue
		}

		if strings.Contains(searchTitle, q) {
			point++
		}

	}

	if en && !areAllEnglish(searchTitle) {
		point -= differentLangsPoint
	}

	return point
}

// areAllEnglish function will check and see if a string contains
// non-english characters or not.
func areAllEnglish(str string) bool {
	for _, s := range str {
		if unicode.IsSymbol(s) || unicode.IsSpace(s) ||
			unicode.IsMark(s) || unicode.IsControl(s) {
			continue
		}

		// some weird characters which Idk?
		if s == weirdChar1 || s == weirdChar2 {
			continue
		}

		if s < endEnglishChars ||
			(s >= startChars1 && s <= endChars1) {
			continue
		} else {
			return false
		}
	}

	return true
}

// getHit returns the best hit.
func getHit(q string, hits []Hits) *Result {
	// make sure that hits array is not nil.
	if hits == nil {
		return nil
	}

	q = strings.ToLower(q)
	en := areAllEnglish(q)

	var best *Result
	var bestPoint int
	var point int
	var current *Result

	for _, h := range hits {
		// since in this function, we are only looking for `song`s,
		// if the type of the hit (or its index) is not "song",
		// the we shouldn't count this loop and will continue our
		// looping.
		if !strings.EqualFold(h.Type, songValue) ||
			!strings.EqualFold(h.Index, songValue) {
			continue
		}

		// make sure our point is zero at the beginning of
		// each loop.
		point = 0
		current = h.Result
		if current == nil {
			// if the current result is nil, continue your loop
			// to find a result which is NOT nil.
			continue
		}

		point += smartFind(q, current.Title, en)
		point += smartFind(q, current.FullTitle, en)
		point += smartFind(q, current.TitleWithFeatured, en)
		point += smartFind(q, current.URL, en)

		if current.LyricsState == completeValue {
			point++
		}

		if point > bestPoint {
			bestPoint = point
			best = current
		}
	}

	if bestPoint <= minConfidence {
		return nil
	}

	return best
}
