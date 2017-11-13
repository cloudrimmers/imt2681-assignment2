package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
	"github.com/gorilla/mux"
)

const root = "/rimbot/"

var err error

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {
	//obtain slacks' webhook
	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.Form

	//convert webhook values into new outgoing message
	response, code := dialogFlow.Query(form.Get("text"))
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}
	test := fmt.Sprintf("Response got:\n%+v", response.Result.Parameters)

	slackTo := struct {
		Text     string `json:"text"`
		Username string `json:"username,omitempty"`
	}{test, "Rimbot"}

	outgoing, err := json.Marshal(slackTo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(outgoing)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rimbot/", Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
