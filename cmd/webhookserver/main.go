package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Arxcis/imt2681-assignment2/lib/handler"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Initializing server....")

	router := mux.NewRouter().StrictSlash(true)
	apiBase := os.Getenv("API_VERSION_PATH")
	port := os.Getenv("PORT")

	//router.HandleFunc("/", handler.HelloWorld).Methods("GET")
	router.HandleFunc(apiBase+"/hook", handler.PostWebhook).Methods("POST")
	router.HandleFunc(apiBase+"/hook", handler.GetWebhookAll).Methods("GET")
	router.HandleFunc(apiBase+"/hook/{id}", handler.GetWebhook).Methods("GET")
	router.HandleFunc(apiBase+"/hook/evaluationtrigger", handler.EvaluationTrigger).Methods("GET")

	router.HandleFunc(apiBase+"/currency/latest", handler.GetLatestCurrency).Methods("GET")
	router.HandleFunc(apiBase+"/currency/average", handler.GetAverageCurrency).Methods("GET")

	log.Println("port: ", port, "apiBase: ", apiBase)
	log.Println(http.ListenAndServe(":"+port, router))
}
