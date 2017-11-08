package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/subosito/gotenv"

	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
)



var APP *App
var err error

func init() {
	gotenv.MustLoad(".env")
	APP = &App{
		FixerioURI:        os.Getenv("FIXERIO_URI"),
		SeedFixerPath:     os.Getenv("SEEDFIXER_PATH"),
		CollectionWebhook: os.Getenv("COLLECTION_FIXER"),
		CollectionFixer:   os.Getenv("COLLECTION_WEBHOOK"),
		MongodbName:	   os.Getenv("MONGODB_NAME"),
		MongodbURI:  	   os.Getenv("MONGODB_URI"),
	}
	mongo := database.Mongo{APP.MongodbName, APP.MongodbURI}
	APP.SeedFixer(&mongo)	 // @note enable/disable
	mongo.EnsureIndex(APP.CollectionFixer, []string{"date"})

	log.Println("Fixer app initialized...")
}

func main() {
	// @doc https://stackoverflow.com/a/35009735
	targetWait := -tool.UntilTomorrow()
	ticker := time.NewTicker(time.Minute)
	mongo := database.Mongo{APP.MongodbName, APP.MongodbURI}
	
	log.Println("T wait  : ", targetWait.String())
	for _ = range ticker.C {
		targetWait += time.Minute

		log.Println("T wait  : ", targetWait.String())

		if targetWait > 0 {
			targetWait = -tool.UntilTomorrow()
			APP.Fixer2Mongo(&mongo)
		}
	}
}
