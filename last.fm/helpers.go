package last_fm

import (
	"encoding/json"
	"fmt"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetRecentTracksByUsername(username string) (res *GetRecentTracks, err error) {
	reqUrl := recentTracksBaseUrl + username + fmt.Sprintf("&extended=1&api_key=%s", config.Data.LastFMKey) + "&format=json"
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
			logging.Warn(err.Error())
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

func GetLastfmTrack(artist string, track string) (res *GetLastFmTrackResponse, err error) {
	reqUrl := getTrackInfoBaseUrl + fmt.Sprintf("&api_key=%s", config.Data.LastFMKey) +
		fmt.Sprintf("&artist=%s&track=%s", url.QueryEscape(artist), url.QueryEscape(track)) + "&format=json"
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// why can't i return anything from here ??
			logging.Warn(err.Error())
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
	reqUrl := userBaseUrl + username + fmt.Sprintf("&api_key=%s", config.Data.LastFMKey) + "&format=json"
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
			logging.Warn(err.Error())
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
			logging.Warn(err.Error())
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
