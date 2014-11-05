package main

import (
	"bufio"
	"fmt"
	"harvester"
	"log"
	"os"
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		p := harvester.NewPayload(url)
		err := p.Download()
		check(err)
		fmt.Println(p.ChildURLs)
	}
	fmt.Println("Done.")
}
