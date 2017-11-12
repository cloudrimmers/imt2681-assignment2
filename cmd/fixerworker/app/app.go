package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
)

// App ...
type App struct {
	FixerioURI          string
	CollectionFixerName string
	Mongo               database.Mongo
	Seedpath            string
}

// Fixer2Mongo ...
func (app *App) Fixer2Mongo() {

	// 1. Connect and request to fixer.io
	resp, err := http.Get(app.FixerioURI)
	if err != nil {
		log.Println("ERROR No connection with fixer.io: "+app.FixerioURI+" ...", err.Error())
		return
	}

	// 2. Decode payload
	payload := &(types.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		log.Println("ERROR Could not decode resp.Body...", err.Error())
		return
	}

	// @TODO 3 Validate incomming data

	// 4. Connect to DB
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	defer app.Mongo.Close()
	if err != nil {
		log.Println("ERROR Database no connection: ", err.Error())
		return
	}

	// 5. Timestamp
	payload.Timestamp = time.Now().String()

	// 6. Dump payload to database
	err = collectionFixer.Insert(payload)
	if err != nil {
		log.Println("ERROR on db.Insert():\n", err.Error())
		return
	}
	log.Println("SUCCESS pulling fixer.io: ", payload)
}

// SeedFixer ...
func (app *App) SeedFixer() {

	// 0. Get seed
	seed := func() []types.FixerIn {
		log.Println("Reading " + app.Seedpath)
		data, err := ioutil.ReadFile(app.Seedpath)
		if err != nil {
			panic(err.Error())
		}
		var fixerin []types.FixerIn
		if err = json.Unmarshal(data, &fixerin); err != nil {
			panic(err.Error())
		}
		log.Println("Done with " + app.Seedpath)
		return fixerin
	}()

	// 1. Open collection
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer app.Mongo.Close()

	// 2. Insert to database
	// cfixer.DropCollection()
	for _, o := range seed {
		if err = collectionFixer.Insert(o); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}
