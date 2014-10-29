package harvester

import (
	"io/ioutil"
	"net/http"
)

func Download(url string, result chan<- []byte, errors chan<- error) {
	resp, err := http.Get(url)
	if err != nil {
		errors <- err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errors <- err
		return
	}
	result <- body
	return
}
