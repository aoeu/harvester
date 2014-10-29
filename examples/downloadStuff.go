package main

import (
	"fmt"
	"harvester"
	"log"
)

func main() {
	results := make(chan []byte)
	errors := make(chan error)
	url := "http://www.google.com"
	go harvester.Download(url, results, errors)
	select {
	case err := <-errors:
		log.Fatal(err)
	case result := <-results:
		fmt.Println(result)
	}
	fmt.Println("Done.")
}
