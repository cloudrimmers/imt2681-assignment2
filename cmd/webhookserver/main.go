package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.MustLoad(".env")
	log.Println("!!! GOTENV !!! ")
}

func main() {
	log.Println("Initializing server....")

	var port = os.Getenv("PORT")
	var apiBase = os.Getenv("API_VERSION_PATH")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", HelloWorld).Methods("GET")
	router.HandleFunc(apiBase+"/webhook", PostWebhook).Methods("POST")
	router.HandleFunc(apiBase+"/webhook", GetWebhookAll).Methods("GET")
	router.HandleFunc(apiBase+"/webhook/{id}", GetWebhook).Methods("GET")
	router.HandleFunc(apiBase+"/webhook/{id}", DeleteWebhook).Methods("DELETE")
	router.HandleFunc(apiBase+"/evaluationtrigger", EvaluationTrigger).Methods("POST")

	router.HandleFunc(apiBase+"/currency/latest", GetLatestCurrency).Methods("POST")
	router.HandleFunc(apiBase+"/currency/average", GetAverageCurrency).Methods("POST")

	log.Println("port: ", port, "apiBase: ", apiBase)
	log.Println(http.ListenAndServe(":"+port, router))
}
