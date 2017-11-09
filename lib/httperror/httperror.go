package httperror

import (
	"log"
	"net/http"
)

// ServiceUnavailable ...
func ServiceUnavailable(w http.ResponseWriter, msg string, err error) {
	log.Println("HTTPERROR  service unavailable | ", msg, " | ", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

// InternalServer ...
func InternalServer(w http.ResponseWriter, msg string, err error) {
	log.Println("HTTPERROR internal server | ", msg, " | ", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// NotFound ...
func NotFound(w http.ResponseWriter, msg string, err error) {
	log.Println("HTTPERROR not found | ", msg, " | ", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, msg string, err error) {
	log.Println("HTTPERROR bad request | ", msg, " | ", err.Error())
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
