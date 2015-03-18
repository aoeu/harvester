package harvester

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
)


func Download(URL string) (body []byte, err error) {
	u, err := url.ParseRequestURI(URL)
	if err != nil {
		// TODO(aoeu): Don't use naked returns.
		return
	}
	resp, err := http.Get(u.String())
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

func ScrapeURLs(body []byte) (URLs map[string]int, err error) {
	URLs = make(map[string]int)
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// TODO: Add base URL, if needed.
		URL, ok := s.Attr("href")
		if ok && ValidateURL(URL) {
			URLs[URL]++
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

// TODO: What is a better name than "Payload" ?

type Payload struct {
	RawResponseBody []byte // Raw response body.
	URL string // The URL used to obtain this payload.
	ChildURLs map[string]int // URLs found within the page and number of times they appeared.
	ScrapeFunction func(body []byte)(map[string]int, error) // An overrideable function for scraping.
}

func NewPayload(URL string) *Payload {
	return &Payload{ 
		RawResponseBody : make([]byte, 0), 
		URL : URL, 
		ChildURLs : make(map[string]int),
		ScrapeFunction : ScrapeURLs,
	}
}

func (p *Payload) Download() (err error) { 
	// TODO: UTF-8 conversion, if necessary, for response body.
	p.RawResponseBody, err = Download(p.URL)
	if err != nil {
		return
	}
	err = p.Scrape()
	return
}

func (p *Payload) Scrape() (err error) {
	p.ChildURLs, err = p.ScrapeFunction(p.RawResponseBody)
	return
}
