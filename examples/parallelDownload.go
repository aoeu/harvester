package main

import (
	"bufio"
	"fmt"
	"harvester"
	"log"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize.
	file, err := os.Open("urls.txt")
	check(err)
	urls := make(chan string)
	results := make(chan []byte)
	errors := make(chan error)
	maxGoRoutines := 3

	// Set up downloader.
	go harvester.ParallelDownload(urls, results, errors, maxGoRoutines)

	// Obtain URLs to download.
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls <- scanner.Text()
		}
		close(urls)
	}()

	// Process results.
	timeoutLength := 3 * time.Second
	timeout := time.After(timeoutLength)
	ever := true
	for ever {
		select {
		case err := <-errors:
			log.Fatal(err)
		case result := <-results:
			fmt.Println(string(result[0:79]))
			timeout = time.After(timeoutLength)
		case <-timeout:
			ever = false
		}
	}

	fmt.Println("Done.")
}
