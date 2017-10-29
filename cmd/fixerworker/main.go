package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/payload"
)

const fixerURL string = "http://api.fixer.io/latest?base=EUR"
const fixerFile string = "../../data/baseEUR.json"
const mongoDbConnect string = "127.0.0.1:27017"

func dumpFromFixerURL() {

	// 1. Connect to DB
	//session, err := mgo.Dial(mongoDbConnect)
	//if err != nil {
	//	log.Fatal("No connection with mongodb @ ", mongoDbConnect)
	//}
	//defer session.Close()

	// 2. Connect and request to
	resp, err := http.Get(fixerURL)
	if err != nil {
		log.Fatal("Wrong contact with: "+fixerURL+" ...", err.Error())
	}

	// 3. Decode payload
	fixerPayload := &(payload.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(fixerPayload)
	if err != nil {
		log.Fatal("Could not decode resp.Body...", err.Error())
	}

	// 4. Dump payload to database
	log.Print(fixerPayload, "\n\n")
}

func main() {

	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	ticker := time.NewTicker(10 * time.Second)
	dumpFromFixerURL()
	for _ = range ticker.C {
		dumpFromFixerURL()
	}
}
