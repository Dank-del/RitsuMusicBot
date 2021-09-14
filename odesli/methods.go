package odesli

import (
	"encoding/json"
	"gitlab.com/Dank-del/lastfm-tgbot/logging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func GetLinks(url string) (r *OdesliResponse, err error) {
	reqUrl := baseUrl + url + country
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

	data := new(OdesliResponse)
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return data, nil

}
