package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/fixerworker/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/subosito/gotenv"

	"github.com/cloudrimmers/imt2681-assignment3/lib/timetool"
)

// READENV read environment from .env file
// @note These 3 bools should probably be command-line arguments - Jonas 13.11.17
const READENV = true

// SEED the database with testdata
const SEED = false

// VERBOSE log more to console
const VERBOSE = true

// APP - configuration data
var APP *app.App
var err error

func init() {
	log.Println("Fixerworker booting up...")

	// 1. Require .env to be present
	if READENV {
		log.Println("Reading .env")
		gotenv.MustLoad(".env")
		log.Println("Done with .env")
	}

	// 2. Initialize the app object
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

	// 3. Default values if empty environment
	if APP.Mongo.URI == "" {
		log.Println("No .env present. Using default values")
		APP.Mongo.URI = "mongodb://localhost"
		APP.Mongo.Name = "test"
	}

	// 4. Ensure index to avoid duplicates
	APP.Mongo.EnsureIndex(APP.CollectionFixerName, []string{"date"})

	// 5. Optional seed and log app object
	if SEED {
		APP.SeedFixerdata()
	}
	if VERBOSE {
		indented, _ := json.MarshalIndent(APP, "", "    ")
		log.Println("App data: ", string(indented))
		log.Println("VERBOSE=false to suppress previous log")
	}

	log.Println("Fixerworker initialized...")
}

func main() {
	// @doc https://stackoverflow.com/a/35009735
	targetWait := -(timetool.UntilTomorrow())
	ticker := time.NewTicker(time.Minute)

	log.Println("T wait  : ", targetWait.String())
	for _ = range ticker.C {
		targetWait += time.Minute
		log.Println("T wait  : ", targetWait.String())

		if targetWait >= 0 {
			targetWait = -(timetool.UntilTomorrow())

			APP.Fixer2Mongo()
		}
	}
}
