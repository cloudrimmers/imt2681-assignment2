package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/fixerworker/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/database"

	"github.com/cloudrimmers/imt2681-assignment3/lib/timetool"
)

// APP - configuration data
var APP *app.App
var err error

func init() {
	// @note These 3 bools should probably be command-line arguments - Jonas 13.11.17
	const SEED = true
	const VERBOSE = true

	// 2. Initialize the app object
	APP = &app.App{
		FixerioURI:          "https://api.fixer.io/latest",
		CollectionFixerName: "fixer",
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
	}

	// 3. Default values if empty environment
	if APP.Mongo.URI == "" {
		APP.Mongo.URI = "mongodb://localhost"
		APP.Mongo.Name = "test"
		log.Println("No .env present. Using default values")
	}

	// 4. Ensure index to avoid duplicates
	APP.Mongo.EnsureIndex(APP.CollectionFixerName, []string{"date"})

	// 5. Optional seed and log app object
	if SEED {
		APP.SeedFixerdata()
		log.Println("Seeded database")
	}
	if VERBOSE {
		indented, _ := json.MarshalIndent(APP, "", "    ")
		log.Println("App data: ", string(indented))
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
			response, err := APP.FixerResponse(APP.FixerioURI)

			if err != nil {
				log.Println("ERROR FixerResponse()")
			} else if err = APP.Fixer2Mongo(response); err != nil {
				log.Println("ERROR Fixer2Mongo()")
			} else {
				log.Println("SUCCESS pulling fixer.io: ", response)
			}
		}
	}
}
