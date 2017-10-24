package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)

	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
