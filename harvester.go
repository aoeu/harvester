package harvester

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
)

func Download(URL string) (body []byte, err error) {
	resp, err := http.Get(URL)
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

func ThrottledDownload(URL string, result chan<- []byte, errors chan<- error,
	throttle chan<- bool) {
	body, err := Download(URL)
	if err != nil {
		errors <- err
		return
	}
	result <- body
	throttle <- true
}

func ParallelDownload(URLs <-chan string, results chan<- []byte,
	errors chan<- error, numGoRoutines int) {
	throttle := make(chan bool, numGoRoutines)
	for i := 0; i < numGoRoutines; i++ {
		throttle <- true
	}
	for URL := range URLs {
		<-throttle
		go ThrottledDownload(URL, results, errors, throttle)
	}
}

func ScrapeURLs(body []byte) (URLs []string, err error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		URL, ok := s.Attr("href")
		if ok && ValidateURL(URL) {
			URLs = append(URLs, URL)
		}
	})
	return
}

func ValidateURL(URL string) bool {
	if URL != "#" {
		return true
	}
	return false
}
