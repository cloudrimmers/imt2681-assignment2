package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const root = "/rimbot/"

var err error

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	//obtain slacks' webhook
	err = r.ParseForm()
	if err != nil {
		status = http.StatusBadRequest
	}
	form := r.Form

	//convert webhook values into new outgoing message

	out := struct {
		Text string `json:"text"`
	}{
		form.Get("text"),
	}

	//Prepare outgoing message
	text, err := json.Marshal(out)
	if err != nil {
		status = http.StatusInternalServerError
	}
	outBody := ioutil.NopCloser(bytes.NewBuffer(text))

	//post message
	http.Post("http://webhook.site/323a95e0-a2f3-4e67-9668-882642f51319", "application/json", outBody)

	w.WriteHeader(status)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rimbot/", Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
