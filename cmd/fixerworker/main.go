package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/fixerworker/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/subosito/gotenv"

	"github.com/cloudrimmers/imt2681-assignment3/lib/tool"
)

// READENV read environment from .env file
const READENV = false

// SEED the database with testdata
const SEED = false

// VERBOSE log more to console
const VERBOSE = false

// APP - configuration data
var APP *app.App
var err error

func init() {
	log.Println("Fixerworker booting up...")

	if READENV {
		log.Println("Reading .env")
		gotenv.MustLoad(".env")
		log.Println("Done with .env")
	}

	// Initialize the app object
	APP = &app.App{
		FixerioURI:          "https://api.fixer.io/latest",
		CollectionFixerName: "fixer",
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
		Seedpath: "./config/fixer.json",
	}

	APP.Mongo.EnsureIndex(APP.CollectionFixerName, []string{"date"})

	if SEED {
		APP.SeedFixer()
	}
	if VERBOSE {
		indented, _ := json.MarshalIndent(APP, "", "    ")
		log.Println(string(indented))
	}

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

		if targetWait >= 0 {
			targetWait = -tool.UntilTomorrow()
			APP.Fixer2Mongo()
		}
	}
}
