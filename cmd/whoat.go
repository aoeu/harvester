package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"harvester"
	"log"
)

type Repos []Repo

type Repo struct {
	Id               int
	Name             string
	Full_name        string // TODO(aoeu): Use annotations and idiomatic names
	Contributors_url string
}

type Contributors []Contributor

type Contributor struct {
	Id                int
	Login             string
	Url               string
	Avatar_url        string
	Organizations_url string
	Repos_url         string
	Contributions     int
}

type R struct {
	Repo
	Contributors
}

var orgName string

func main() {
	flag.StringVar(&orgName, "organization name", "hackerschool", "A github organization name.")
	flag.Parse()

	reposURI := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)
	body, err := harvester.Download(reposURI)
	if err != nil {
		log.Fatal(err)
	}

	var repos Repos
	err = json.Unmarshal(body, &repos)
	if err != nil {
		log.Fatal(err)
	}

	repoBundles := make([]R, len(repos))
	for i, repo := range repos {
		repoBundles[i].Repo = repo

		body, err := harvester.Download(repo.Contributors_url)
		if err != nil {
			log.Fatal(err)
		}
		var c Contributors
		err = json.Unmarshal(body, &c)
		if err != nil {
			log.Fatal(err)
		}
		repoBundles[i].Contributors = c
	}

	fmt.Printf("%+v\n", repoBundles)
}
