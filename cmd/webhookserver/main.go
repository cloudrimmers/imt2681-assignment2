package main

import (
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
	gotenv.MustLoad(".env")
	APP = &App{
		CollectionWebhook: os.Getenv("COLLECTION_FIXER"),
		CollectionFixer:   os.Getenv("COLLECTION_WEBHOOK"),
		MongodbName:       os.Getenv("MONGODB_NAME"),
		MongodbURI:        os.Getenv("MONGODB_URI"),
		Mongo:             database.Mongo{Name: APP.MongodbName, URI: APP.MongodbURI, Session: nil},
	}

	log.Println("Initialized server...")
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
