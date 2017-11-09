package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/invoke"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
)

// App ...
type App struct {
	FixerioURI        string
	CollectionWebhook string
	CollectionFixer   string
	Mongo             database.Mongo
	Seed              []types.FixerIn
}

// Fixer2Mongo ...
func (app *App) Fixer2Mongo() {

	// 1. Connect and request to fixer.io
	resp, err := http.Get(app.FixerioURI)
	if err != nil {
		log.Println("ERROR No connection with fixer.io: "+APP.FixerioURI+" ...", err.Error())
		return
	}

	// 2. Decode payload
	payload := &(types.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		log.Println("ERROR Could not decode resp.Body...", err.Error())
		return
	}

	// 3. Connect to DB
	dbsession, err := app.Mongo.Open()
	defer app.Mongo.Close()
	if err != nil {
		log.Println("ERROR Database no connection: ", err.Error())
		return
	}

	payload.Timestamp = time.Now().String()

	// 5. Dump payload to database
	err = dbsession.C(APP.CollectionFixer).Insert(payload)
	if err != nil {
		log.Println("ERROR on db.Insert():\n", err.Error())
		return
	}
	log.Println("SUCCESS pulling fixer.io: ", payload)

	// 6. Fire webhooks
	client := &http.Client{}
	invoke.Webhooks(client, dbsession.C(APP.CollectionWebhook), dbsession.C(APP.CollectionFixer))
}

// SeedFixer ...
func (app *App) SeedFixer() {

	// 1. Open collection
	collection, err := app.Mongo.OpenC(app.CollectionFixer)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer app.Mongo.Close()

	// 2. Insert to database
	// cfixer.DropCollection()
	for _, o := range app.Seed {
		if err = collection.Insert(o); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}
