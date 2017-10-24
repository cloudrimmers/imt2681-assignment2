package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const dataSource string = "http://api.fixer.io/latest?base=EUR"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)
	res, err := http.Get(dataSource)
	if err != nil {
		log.Fatal("Wrong contact with: http://api.fixer.io/latest?base=EUR...", err.Error())
	}
	marshalled, __ := json.Unmarshal(res.Body)
	if err != nil {
		panic(err.Error())
	}

	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
