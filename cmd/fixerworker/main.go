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

var config *types.WebConfig = (&types.WebConfig{}).Load()

func main() {

	database.EnsureFixerIndex(config.CollectionFixer) // @note you may only do this when needed
	// database.SeedFixer(config.CollectionFixer)     // @note only do this when needed
	// @doc https://stackoverflow.com/a/35009735
	log.Println("Initializing ticker...")

	ticker := time.NewTicker(time.Minute)
	targetWait := tool.UntilTomorrow()
	var currentWait time.Duration
	log.Println("Current wait : ", currentWait.String())
	log.Println("Target wait  : ", targetWait.String())
	for _ = range ticker.C {
		currentWait += time.Minute

		log.Println("Current wait : ", currentWait.String())
		log.Println("Target wait  : ", targetWait.String())

		if currentWait >= targetWait {
			targetWait = tool.UntilTomorrow()
			currentWait = 0
			fixer2mongo(os.Getenv("FIXERIO_URI"))
		}
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

	payload.Timestamp = time.Now().String()

	// 5. Dump payload to database
	log.Println(config.CollectionFixer)
	err = db.C(config.CollectionFixer).Insert(payload)
	if err != nil {
		log.Println("Error on db.Insert():\n", err.Error())
		return
	}
	log.Println("Successfull grab of fixer.io: ", payload.Datestamp)

	// 6. Fire webhooks
	tool.FireWebhooks(db.C(config.CollectionWebhook), db.C(config.CollectionFixer))
}
