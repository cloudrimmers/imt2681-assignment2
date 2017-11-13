package main

import (
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/rimbot/app"
	"github.com/gorilla/mux"
)

const root = "/rimbot/"

func main() {
	r := mux.NewRouter()
	r.HandleFunc(root, app.Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
