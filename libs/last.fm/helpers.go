package last_fm

import (
	"encoding/json"
	"fmt"
	"gitlab.com/Dank-del/lastfm-tgbot/core/config"
	"gitlab.com/Dank-del/lastfm-tgbot/core/logging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const APILimit = 1000

func GetRecentTracksByUsername(username string, limit int) (res *GetRecentTracks, err error) {
	if limit > APILimit {
		limit = APILimit
	}
	reqUrl := recentTracksBaseUrl + url.QueryEscape(username) +
		fmt.Sprintf("&extended=1&api_key=%s", config.Data.LastFMKey) +
		"&limit=" + strconv.Itoa(limit) +
		"&format=json"
	//fmt.Println(reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.SUGARED.Warn(err.Error())
		}
	}(resp.Body)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data := new(GetRecentTracks)
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data, nil

}

func GetTopTracks(username string) (res *TopTrackResponse, err error) {
	reqUrl := getTopTracksForUserBaseUrl + username + fmt.Sprintf("&api_key=%s", config.Data.LastFMKey)
	// print(reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.SUGARED.Warn(err.Error())
		}
	}(resp.Body)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data := new(TopTrackResponse)
	err = json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetLastfmTrack(artist string, track string, username string) (res *GetLastFmTrackResponse, err error) {
	reqUrl := getTrackInfoBaseUrl + fmt.Sprintf("&api_key=%s", config.Data.LastFMKey) +
		fmt.Sprintf("&artist=%s&track=%s&username=%s", url.QueryEscape(artist), url.QueryEscape(track), url.QueryEscape(username)) + "&format=json"
	resp, err := http.Get(reqUrl)
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.SUGARED.Warn(err.Error())
		}
	}(resp.Body)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data := new(GetLastFmTrackResponse)
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data, nil
}

func GetLastFMUser(username string) (res *LastFMUser, err error) {
	reqUrl := userBaseUrl + url.QueryEscape(username) +
		fmt.Sprintf("&api_key=%s", config.Data.LastFMKey) + "&format=json"
	//logging.Info(reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.SUGARED.Warn(err.Error())
		}
	}(resp.Body)

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data := new(LastFMUser)
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data, nil
}

func GetTopArtists() (res *GetTopArtistsResponse, err error) {
	reqUrl := getTopArtistsBaseURl + fmt.Sprintf("&api_key=%s", config.Data.LastFMKey) + "&format=json"
	fmt.Println(reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.SUGARED.Warn(err.Error())
		}
	}(resp.Body)
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data := new(GetTopArtistsResponse)
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data, nil
}
