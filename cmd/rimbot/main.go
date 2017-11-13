package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/cmd/rimbot/app"
	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
	"github.com/cloudrimmers/imt2681-assignment3/lib/validate"
	"github.com/gorilla/mux"
)

const root = "/rimbot/"

var err error

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	weReSorry := fmt.Sprint("Sorry, I missed that. Maybe something was vague with what you said?\n",
		"Try using capital letters like this: 'USD', 'GBP'. And numbers like this: '131.5'")

	//obtain slacks' webhook
	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.Form

	slackTo := struct {
		Text     string `json:"text"`
		Username string `json:"username,omitempty"`
	}{}

	//convert webhook values into new outgoing message
	base, target, amount, code := dialogFlow.Query(form.Get("text"))

	if code != http.StatusOK && code != http.StatusPartialContent {
		w.WriteHeader(code)
		return
	} else if code == http.StatusPartialContent { //If Unmarshal fails in Query(). Meaning Clara got confused.
		slackTo.Text = weReSorry
	} else { //If everything got parsed correctly.
		errBase := validate.Currency(base)
		errTarget := validate.Currency(target)

		if errBase == nil && errTarget == nil && amount >= 0 { //If valid input for currencyservice.

			slackTo.Text = fmt.Sprintf("Response got:\n base: %v\ttarget: %v\tamount: %v", base, target, amount) //Temporary outprint
			//Here goes communication with Currencyservice.

		} else { //If invalid input for currencyservice.
			slackTo.Text = weReSorry
		}
	}

	outgoing, err := json.Marshal(slackTo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(outgoing)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc(root, app.Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
