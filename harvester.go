package harvester

import (
	"io/ioutil"
	"net/http"
)

func Download(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func ThrottledDownload(url string, result chan<- []byte, errors chan<- error,
	throttle chan<- bool) {
	body, err := Download(url)
	if err != nil {
		errors <- err
		return
	}
	result <- body
	throttle <- true
}

func ParallelDownload(urls <-chan string, results chan<- []byte,
	errors chan<- error, numGoRoutines int) {
	throttle := make(chan bool, numGoRoutines)
	for i := 0; i < numGoRoutines; i++ {
		throttle <- true
	}
	for url := range urls {
		<-throttle
		go ThrottledDownload(url, results, errors, throttle)
	}
}
