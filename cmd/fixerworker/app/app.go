package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"github.com/cloudrimmers/imt2681-assignment3/lib/validate"
)

// App ...
type App struct {
	FixerioURI          string
	CollectionFixerName string
	Mongo               database.Mongo
}

var err error

// FixerResponse ...
func (app *App) FixerResponse(uri string) (*types.FixerIn, error) {
	// 1. Connect and request to fixer.io
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	// 2. Decode payload
	//rbody := new(types.FixerIn)
	rbody := types.FixerIn{}

	log.Println("!! -", &rbody, resp.Body, "-- !!")

	err = json.NewDecoder(resp.Body).Decode(&rbody)
	
	log.Println("!! --- !!")
	if err != nil {
		return nil, err
	}
	log.Println("!! --- !!")

	return &rbody, nil
}

// Fixer2Mongo ...
func (app *App) Fixer2Mongo(response *types.FixerIn) error {

	// @TODO 3 Validate incomming data
	if err = validate.Currency(response.Base); err != nil {
		return err
	}

	// 4. Connect to DB
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		return err
	}
	defer app.Mongo.Close()

	// 5. Timestamp
	response.Timestamp = time.Now().String()

	// 6. Dump payload to database
	err = collectionFixer.Insert(response)
	if err != nil {
		return err
	}

	return nil
}

// SeedFixerdata ...
// @TODO This is a duplicate of an identical function in fixerworker/
//   	  Maybe this should be put in a library - JSolsvik 15.11.17
func (app *App) SeedFixerdata() {

	// 0. Get seed
	seed := types.FixerSeed

	// 1. Open collection
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer app.Mongo.Close()
	collectionFixer.DropCollection()

	// 2. Insert to database
	// cfixer.DropCollection()
	for _, o := range seed {
		if err = collectionFixer.Insert(o); err != nil {
			log.Println("Unable to db.Insert seed", err.Error())
		}
	}
}
