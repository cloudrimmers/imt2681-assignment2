package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arxcis/imt2681-assignment2/lib/handler"
	"github.com/Arxcis/imt2681-assignment2/lib/payload"
	"github.com/gorilla/mux"
)

func runTicker(done chan string, tickerCallback func(*payload.Ticker)) {
	log.Println("Initializing ticker....")

	// @doc https://stackoverflow.com/a/35009735
	// ticker1 := time.NewTicker(24 * time.Hour)  // This should be under production
	ticker1 := time.NewTicker(60 * time.Second)

	for _ = range ticker1.C {

		resp, err := http.Get(dataSource)
		if err != nil {
			log.Fatal("Wrong contact with: http://api.fixer.io/latest?base=EUR ...", err.Error())
		}

		ticker := &(payload.Ticker{})
		err = json.NewDecoder(resp.Body).Decode(ticker)

		if err != nil {
			log.Fatal("Could not decode resp.Body...", err.Error())
		}
		tickerCallback(ticker)
	}
	done <- "Done"
}

func runServer(done chan string) {
	log.Println("Initializing server....")

	router := mux.NewRouter().StrictSlash(true)
	handler.InitHandlers(router)
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))

	done <- "Done"
}

func main() {

	done := make(chan string)

	go runTicker(done, func(payload *payload.Ticker) { log.Print(payload, "\n\n") })
	go runServer(done)

	<-done
}
