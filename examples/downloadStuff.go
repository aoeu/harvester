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
	// This won't return since we didn't start it with the go keyword.
	// In other words, the unbuffered channel will not be written to
	// if nothing is yet reading from the unbuffered channel.
	harvester.Download(url, results, errors)
	select {
	case err := <-errors:
		log.Fatal(err)
	case result := <-results:
		fmt.Println(result)
	}
	fmt.Println("Done.")
}
