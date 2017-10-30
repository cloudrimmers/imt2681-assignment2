package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	mgo "gopkg.in/mgo.v2"
)

func fixer2mongo(
	mongoURI string,
	fixerURI string,
	mongoDB string,
	mongoC string) {

	// 1. Connect and request to fixer.io
	resp, err := http.Get(fixerURI)
	if err != nil {
		log.Println("Wrong contact with: "+fixerURI+" ...", err.Error())
		return
	}

	// 2. Decode payload
	payload := &(mytypes.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		log.Println("Could not decode resp.Body...", err.Error())
		return
	}

	// 3. Connect to DB
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.Println("No connection with mongodb @ ", mongoURI, err.Error())
		return
	}
	defer session.Close()

	// 4. Generate datestamp
	now := time.Now()
	payload.Datestamp = fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

	// 5. Dump payload to database
	err = session.DB(mongoDB).C(mongoC).Insert(payload)
	if err != nil {
		log.Println("Error on db.Insert():\n", err.Error())
		return
	}

	log.Print("Tick success: ", payload.Datestamp)
}

func main() {

	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	for {
		ticker := time.NewTicker(tool.UntilTomorrow())
		<-ticker.C // Wait
		ticker.Stop()

		fixer2mongo(
			os.Getenv("MONGODB_URI"),
			os.Getenv("FIXERIO_URI"),
			os.Getenv("MONGODB_NAME"),
			os.Getenv("MONGODB_COLLECTION"))
	}
}
