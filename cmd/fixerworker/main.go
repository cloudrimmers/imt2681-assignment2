package main

import (
	"log"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/subosito/gotenv"

	"github.com/Arxcis/imt2681-assignment2/lib/tool"
)

// APP - global state pbject
var APP *App
var err error

func init() {
	gotenv.MustLoad(".env")

	APP = &App{
		FixerioURI:        os.Getenv("FIXERIO_URI"),
		SeedFixerPath:     "./config/seedfixer.json",
		CollectionWebhook: os.Getenv("COLLECTION_FIXER"),
		CollectionFixer:   os.Getenv("COLLECTION_WEBHOOK"),
		MongodbName:       os.Getenv("MONGODB_NAME"),
		MongodbURI:        os.Getenv("MONGODB_URI"),
		Mongo:             database.Mongo{Name: APP.MongodbName, URI: APP.MongodbURI, Session: nil},
	}
	APP.SeedFixer() // @note enable/disable
	APP.Mongo.EnsureIndex(APP.CollectionFixer, []string{"date"})

	log.Println("Fixer app initialized...")
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
