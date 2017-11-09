package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/cmd/fixerworker/app"
	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/subosito/gotenv"

	"github.com/Arxcis/imt2681-assignment2/lib/tool"
)

// APP - global state pbject
var APP *app.App
var err error

func init() {
	log.Println("Fixerworker booting up...")

	log.Println("Reading .env")
//	gotenv.MustLoad(".env")
	log.Println("Done with .env")

	configpath := "./config/seedfixer.json"
	APP = &app.App{
		FixerioURI:        os.Getenv("FIXERIO_URI"),
		CollectionWebhook: os.Getenv("COLLECTION_WEBHOOK"),
		CollectionFixer:   os.Getenv("COLLECTION_FIXER"),
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
		Seed: func() []types.FixerIn {
			log.Println("Reading " + configpath)
			data, err := ioutil.ReadFile(configpath)
			if err != nil {
				panic(err.Error())
			}
			var fixerin []types.FixerIn
			if err = json.Unmarshal(data, &fixerin); err != nil {
				panic(err.Error())
			}
			log.Println("Done with " + configpath)
			return fixerin
		}(),
	}
	//APP.SeedFixer() // @note enable/disable
	APP.Mongo.EnsureIndex(APP.CollectionFixer, []string{"date"})

	//@verbose
	//	indented, _ := json.MarshalIndent(APP, "", "    ")
	//	log.Println(string(indented))
	log.Println("Fixerworker initialized...")
}

func main() {
	// @doc https://stackoverflow.com/a/35009735
	targetWait := -tool.UntilTomorrow()
	ticker := time.NewTicker(time.Minute)

	log.Println("T wait  : ", targetWait.String())
	for _ = range ticker.C {
		targetWait += time.Minute
		log.Println("T wait  : ", targetWait.String())

		if targetWait > 0 {
			targetWait = -tool.UntilTomorrow()
			APP.Fixer2Mongo()
		}
	}
}
