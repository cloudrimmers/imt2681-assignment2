package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/payload"
	mgo "gopkg.in/mgo.v2"
)

const fixerURL string = "http://api.fixer.io/latest?base=EUR"
const fixerFile string = "../../data/baseEUR.json"
const mongoDbConnect string = "127.0.0.1:27017"

func dumpFromFixerURL() {

	// 1. Connect and request to fixer.io
	resp, err := http.Get(fixerURL)
	if err != nil {
		log.Fatal("Wrong contact with: "+fixerURL+" ...", err.Error())
		return
	}

	// 2. Decode payload
	fixerPayload := &(payload.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(fixerPayload)
	if err != nil {
		log.Fatal("Could not decode resp.Body...", err.Error())
		return
	}

	// 3. Connect to DB
	session, err := mgo.Dial(mongoDbConnect)
	if err != nil {
		log.Fatal("No connection with mongodb @ ", mongoDbConnect, err.Error())
		return
	}
	defer session.Close()

	// 4. Dump payload to database
	err = session.DB("test").C("tick").Insert(fixerPayload)
	if err != nil {
		log.Fatal("Error on session.DB().C('tick').Insert(<Payload>)", fixerPayload, err.Error())
	}

	log.Print("Tick success: ", fixerPayload)
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
