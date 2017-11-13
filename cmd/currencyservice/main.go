package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/currencyservice/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/database"

	"github.com/gorilla/mux"
)

// APP - global state pbject
var APP *app.App

func init() {
	const VERBOSE = true
	const configpath = "./config/currency.json"

	APP = &app.App{
		CollectionFixerName: "fixer",
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
		Currency: func() []string {
			log.Println("Reading " + configpath)
			data, err := ioutil.ReadFile(configpath)
			if err != nil {
				panic(err.Error())
			}
			var currency []string
			if err = json.Unmarshal(data, &currency); err != nil {
				panic(err.Error())
			}
			log.Println("Done with " + configpath)
			return currency
		}(),
	}

	// 3. Default values if empty environment
	if APP.Mongo.URI == "" {
		log.Println("No .env present. Using default values")
		APP.Mongo.URI = "mongodb://localhost"
		APP.Mongo.Name = "test"
		APP.Port = "5000"
	}

	if VERBOSE {
		indented, _ := json.MarshalIndent(APP, "", "    ")
		log.Println("App data: ", string(indented))
	}
	log.Println("Currencyservice initialized...")
}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/currency/latest", APP.GetLatestCurrency).Methods("GET")
	log.Println(http.ListenAndServe(":"+APP.Port, router))
}
