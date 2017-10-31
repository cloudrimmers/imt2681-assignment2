package main

import (
	"log"
	"net/http"
)

func serviceUnavailable(w http.ResponseWriter, err error) {
	log.Println("No database connection: ", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println("Collection.Insert() error", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func notFound(w http.ResponseWriter, err error) {
	log.Println("Collection.Find() not found", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println("Http bad request", err.Error())
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
