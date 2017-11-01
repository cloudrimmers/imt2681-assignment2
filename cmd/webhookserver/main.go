package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Initializing server....")

	router := mux.NewRouter().StrictSlash(true)
	apiBase := os.Getenv("API_VERSION_PATH")
	port := os.Getenv("PORT")

	//router.HandleFunc("/", HelloWorld).Methods("GET")
	router.HandleFunc(apiBase+"/hook", PostWebhook).Methods("POST")
	router.HandleFunc(apiBase+"/hook", GetWebhookAll).Methods("GET")
	router.HandleFunc(apiBase+"/hook/{id}", GetWebhook).Methods("GET")
	router.HandleFunc(apiBase+"/hook/evaluationtrigger", EvaluationTrigger).Methods("POST")

	router.HandleFunc(apiBase+"/currency/latest", GetLatestCurrency).Methods("POST")
	router.HandleFunc(apiBase+"/currency/average", GetAverageCurrency).Methods("POST")

	log.Println("port: ", port, "apiBase: ", apiBase)
	log.Println(http.ListenAndServe(":"+port, router))
}
