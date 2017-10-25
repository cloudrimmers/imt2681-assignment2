package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/payload"
)

func dumpFromURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Wrong contact with: "+url+" ...", err.Error())
	}

	fixerPayload := &(payload.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(fixerPayload)

	if err != nil {
		log.Fatal("Could not decode resp.Body...", err.Error())
	}
	log.Print(fixerPayload, "\n\n")
}

func main() {

	const fixerURL string = "http://api.fixer.io/latest?base=EUR"
	const fixerFile string = "../../data/baseEUR.json"

	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	ticker := time.NewTicker(60 * time.Minute)

	dumpFromURL(fixerURL)
	for _ = range ticker.C {
		dumpFromURL(fixerURL)
	}
}
