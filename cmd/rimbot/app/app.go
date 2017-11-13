package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
)

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
	base, target, amount, code := dialogFlow.Query(form.Get("text"))
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}
	test := fmt.Sprintf("Response got:\n%v\t%v\t%v", base, target, amount)

	//Here goes validation of base, target, and amount.

	//Here goes communication with Currencyservice.

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
