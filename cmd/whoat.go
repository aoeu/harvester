package main 

import (
	"fmt"
	"flag"
	"log"
	"net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Repos []Repo

type Repo struct {
	Id int
	Name string
	Full_name string
	Contributors_url string
}


var orgName string

func main() { 
	flag.StringVar(&orgName, "organization name", "hackerschool", "A github organization name.")
	flag.Parse()

	s := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)
	reposURL, err := url.ParseRequestURI(s)
	if err != nil {
		log.Fatal(err)
	}
	
	resp, err := http.Get(reposURL.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	var repos Repos
	err = json.Unmarshal(body, &repos)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", repos)
}
