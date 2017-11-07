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

	log.Println("Initializing ticker...")

	database.EnsureFixerIndex() // @note you may only do this when needed
	// database.SeedFixer(config.CollectionFixer)     // @note only do this when needed
	// @doc https://stackoverflow.com/a/35009735

	targetWait := -tool.UntilTomorrow()
	log.Println("T wait  : ", targetWait.String())

	//	fixer2mongo(os.Getenv("FIXERIO_URI"))

	ticker := time.NewTicker(time.Minute)
	for _ = range ticker.C {
		targetWait += time.Minute

		log.Println("T wait  : ", targetWait.String())

		if targetWait > 0 {
			targetWait = -tool.UntilTomorrow()
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

	}
	log.Println("Successfull grab of fixer.io: ", payload)

	// 6. Fire webhooks
	client := &http.Client{}
	tool.InvokeWebhooks(client, db.C(os.Getenv("COLLECTION_WEBHOOK")), db.C(os.Getenv("COLLECTION_FIXER")))
}
