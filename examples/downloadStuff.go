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
	results := make(chan []byte)
	errors := make(chan error)

	file, err := os.Open("urls.txt")
	check(err)

	throttleSize := 3
	throttle := make(chan bool, throttleSize)
	for i := 0; i < throttleSize; i++ {
		throttle <- true
	}
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			<-throttle
			go harvester.ThrottledDownload(url, results, errors, throttle)
		}
	}()

	timeoutLength := 3 * time.Second
	timeout := time.After(timeoutLength)
	ever := true
	for ever {
		select {
		case err := <-errors:
			log.Fatal(err)
		case result := <-results:
			fmt.Println(result)
			timeout = time.After(timeoutLength)
		case <-timeout:
			ever = false
		}
	}

	fmt.Println("Done.")
}
