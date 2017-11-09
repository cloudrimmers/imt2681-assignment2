package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Arxcis/imt2681-assignment2/lib/database"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

// APP - global state pbject
var APP *App

func init() {
	log.Println("Webhookserver starting up...")

	log.Println("Reading .env")
	gotenv.MustLoad(".env")
	log.Println("Done with .env")

	configpath := "./config/currency.json"
	APP = &App{
		CollectionWebhook: os.Getenv("COLLECTION_FIXER"),
		CollectionFixer:   os.Getenv("COLLECTION_WEBHOOK"),
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
	indented, _ := json.MarshalIndent(APP, "", "    ")
	log.Println(string(indented))
	log.Println("Webhookserver initialized...")
}

func main() {

	var port = os.Getenv("PORT")
	var apiBase = os.Getenv("API_VERSION_PATH")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", APP.HelloWorld).Methods("GET")
	router.HandleFunc(apiBase+"/webhook", APP.PostWebhook).Methods("POST")
	router.HandleFunc(apiBase+"/webhook", APP.GetWebhookAll).Methods("GET")
	router.HandleFunc(apiBase+"/webhook/{id}", APP.GetWebhook).Methods("GET")
	router.HandleFunc(apiBase+"/webhook/{id}", APP.DeleteWebhook).Methods("DELETE")
	router.HandleFunc(apiBase+"/webhook/evaluationtrigger", APP.EvaluationTrigger).Methods("POST")

	router.HandleFunc(apiBase+"/currency/latest", APP.GetLatestCurrency).Methods("POST")
	router.HandleFunc(apiBase+"/currency/average", APP.GetAverageCurrency).Methods("POST")

	log.Println("port: ", port, "apiBase: ", apiBase)
	log.Println(http.ListenAndServe(":"+port, router))
}
