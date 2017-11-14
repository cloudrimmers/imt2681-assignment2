package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/rimbot/app"
	"github.com/gorilla/mux"
)

const root = "/rimbot/"

var err error

func main() {
	r := mux.NewRouter()
	r.HandleFunc(root, app.Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)
	log.Println("RIMPORT: ", os.Getenv("RIMPORT"))
	http.ListenAndServe(":"+os.Getenv("RIMPORT"), nil)
}
