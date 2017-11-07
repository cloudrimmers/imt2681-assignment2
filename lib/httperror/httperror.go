package httperror

import (
	"log"
	"net/http"
)

// ServiceUnavailable ...
func ServiceUnavailable(w http.ResponseWriter, err error) {
	log.Println("No database connection: ", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

// InternalServer ...
func InternalServer(w http.ResponseWriter, err error) {
	log.Println("Collection.Insert() error", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// NotFound ...
func NotFound(w http.ResponseWriter, err error) {
	log.Println("Collection.Find() not found", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, err error) {
	log.Println("Http bad request", err.Error())
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
