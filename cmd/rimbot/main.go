package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const root = "/rimbot/"
const dialogFlowRoot = "https://api.dialogflow.com/v1/query" //NOTE: protocol number is "required", consider adding it

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

	//TODO generate sessionId every time.

	out := struct {
		Query     string `json:"query"`
		SessionID int    `json:"sessionId"`
	}{form.Get("text"), rand.Intn(10000)} //Generate random SessionID.

	fmt.Printf("%+v\n", out) //Print the body that will be sent to DialigFlow.

	//Prepare outgoing message
	text, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outBody := ioutil.NopCloser(bytes.NewBuffer(text)) //Create new body with text to send dialogflow.

	//post message
	req, err := http.NewRequest(http.MethodPost, dialogFlowRoot, outBody) //Creates new request with body.

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN")) //Add authorization token to head. Identifies agent in dialogflow.
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req) //Execute request.
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency) //NOTE: is this right?
		return
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(respBody)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rimbot/", Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
