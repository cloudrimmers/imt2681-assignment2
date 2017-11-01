package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/database"

	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
)

func main() {

	database.EnsureFixerIndex()

	// @doc https://stackoverflow.com/a/35009735
	fixer2mongo(os.Getenv("FIXERIO_URI")) // Seed database
	log.Println("Initializing ticker...")

	for {
		ticker := time.NewTicker(tool.UntilTomorrow())
		<-ticker.C // Wait
		ticker.Stop()

		fixer2mongo(os.Getenv("FIXERIO_URI"))
	}
}

func fixer2mongo(fixerURI string) {

	// 1. Connect and request to fixer.io
	resp, err := http.Get(fixerURI)
	if err != nil {
		log.Println("No connection with fixer.io: "+fixerURI+" ...", err.Error())
		return
	}

	// 2. Decode payload
	payload := &(types.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		log.Println("Could not decode resp.Body...", err.Error())
		return
	}

	// 3. Connect to DB
	db, err := database.Open()
	if err != nil {
		log.Println("Database no connection: ", err.Error())
	}
	defer database.Close()

	// 4. Generate datestamp
	payload.Datestamp = tool.Todaystamp()

	// 5. Dump payload to database
	err = db.C("tick").Insert(payload)
	if err != nil {
		log.Println("Error on db.Insert():\n", err.Error())
		return
	}

	log.Print("Tick success: ", payload.Datestamp)
}
