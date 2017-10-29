package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/payload"
	mgo "gopkg.in/mgo.v2"
)

const fixerURL string = "http://api.fixer.io/latest?base=EUR"
const mongoDbName string = "test"
const mongoDbCollection string = "tick"
const intervall time.Duration = 24 * time.Hour

var mongoDbConnect string = os.Getenv("MONGODB_URI")

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
	err = session.DB(mongoDbName).C(mongoDbCollection).Insert(fixerPayload)
	if err != nil {
		log.Fatal("Error on session.DB(", mongoDbName, ").C(", mongoDbCollection, ").Insert(<Payload>)", fixerPayload, err.Error())
	}

	log.Print("Tick success: ", fixerPayload)
}

func main() {

	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	ticker := time.NewTicker(intervall)
	dumpFromFixerURL()
	for _ = range ticker.C {
		dumpFromFixerURL()
	}
}
