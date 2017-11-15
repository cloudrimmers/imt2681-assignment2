package httperror

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

// ServiceUnavailable ...
func ServiceUnavailable(w http.ResponseWriter, msg string, err error) {
	_, fn, line, _ := runtime.Caller(1)
	_, fname := filepath.Split(fn)

	log.Println("HTTPERROR|", http.StatusText(http.StatusServiceUnavailable), "| ", fname, line, msg, "|", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

// InternalServer ...
func InternalServer(w http.ResponseWriter, msg string, err error) {
	_, fn, line, _ := runtime.Caller(1)
	_, fname := filepath.Split(fn)

	log.Println("HTTPERROR", http.StatusText(http.StatusInternalServerError), "|", fname, line, msg, "|", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// NotFound ...
func NotFound(w http.ResponseWriter, msg string, err error) {
	_, fn, line, _ := runtime.Caller(1)
	_, fname := filepath.Split(fn)

	log.Println("HTTPERROR", http.StatusText(http.StatusNotFound), "|", fname, line, msg, "|", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, msg string, err error) {
	_, fn, line, _ := runtime.Caller(1)
	_, fname := filepath.Split(fn)

	log.Println("HTTPERROR", http.StatusText(http.StatusBadRequest), "|", fname, line, msg, "|", err.Error())
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
