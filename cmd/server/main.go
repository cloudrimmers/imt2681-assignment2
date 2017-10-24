package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/payload"
	"github.com/gorilla/mux"
)

const dataSource string = "http://api.fixer.io/latest?base=EUR"

func initializeTicker(messages chan string, tickerCallback func(*payload.Ticker)) {
	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	// ticker1 := time.NewTicker(24 * time.Hour)  // This should be under production
	ticker1 := time.NewTicker(60 * time.Second)

	for _ = range ticker1.C {

		resp, err := http.Get(dataSource)
		if err != nil {
			log.Fatal("Wrong contact with: http://api.fixer.io/latest?base=EUR ...", err.Error())
		}

		ticker := payload.Ticker{}
		err = json.NewDecoder(resp.Body).Decode(&ticker)

		if err != nil {
			log.Fatal("Could not decode resp.Body...", err.Error())
		}
		tickerCallback(&ticker)
	}
	messages <- "Done"
}

func initializeServer(messages chan string) {
	router := mux.NewRouter().StrictSlash(true)
	// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)

	// Method  Route   About
	// GET /api/v1/subscription/   list subscriptions
	// POST    /api/v1/subscription/   create a subscription
	// GET /api/v1/subscription/:id/   get a subscription
	// PUT /api/v1/subscription/:id/   update a subscription
	// DELETE  /api/v1/subscription/:id/   delete a subscription

	log.Println("Initializing server....")
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))

	messages <- "Done"
}

func main() {

	messages := make(chan string)

	go initializeTicker(messages, func(payload *payload.Ticker) { log.Print(payload, "\n\n") })
	go initializeServer(messages)

	<-messages
}
