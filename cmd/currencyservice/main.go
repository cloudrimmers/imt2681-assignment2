package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/currencyservice/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/environment"

	"github.com/gorilla/mux"
)

// APP - global state pbject
var APP *app.App
var err error

func init() {
	// 1. Load environment
	if err = environment.Load(os.Args); err != nil {
		panic(err.Error())
	}

	APP = &app.App{
		Port:                os.Getenv("PORT"),
		CollectionFixerName: "fixer",
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
	}

	indented, _ := json.MarshalIndent(APP, "", "    ")
	log.Println("App data: ", string(indented))

	log.Println("Currencyservice initialized...")
}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/currency/latest/", APP.GetLatestCurrency).Methods("POST")
	log.Println("PORT: ", os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+APP.Port, router))
}
