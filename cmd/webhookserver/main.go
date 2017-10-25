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
	handler.InitHandlers(router)
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
