package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Initializing server....")

	port := os.Getenv("PORT")
	apiBase := os.Getenv("API_VERSION_PATH")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", HelloWorld).Methods("GET")
	router.HandleFunc(apiBase+"/hook", PostWebhook).Methods("POST")
	router.HandleFunc(apiBase+"/hook", GetWebhookAll).Methods("GET")
	router.HandleFunc(apiBase+"/hook/{id}", GetWebhook).Methods("GET")
	router.HandleFunc(apiBase+"/hook/{id}", DeleteWebhook).Methods("DELETE")
	router.HandleFunc(apiBase+"/hook/evaluationtrigger", EvaluationTrigger).Methods("POST")

	router.HandleFunc(apiBase+"/currency/latest", GetLatestCurrency).Methods("POST")
	router.HandleFunc(apiBase+"/currency/average", GetAverageCurrency).Methods("POST")

	log.Println("port: ", port, "apiBase: ", apiBase)
	log.Println(http.ListenAndServe(":"+port, router))
}
