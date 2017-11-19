package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
	"github.com/cloudrimmers/imt2681-assignment3/lib/validate"
)

var err error

const slackUserError = "Sorry, I missed that. Maybe something was vague with what you said? Try using capital letters like this: 'USD', 'GBP'. And numbers like this: '131.5'"
const BotDefaultName = "Rimbot"

func MessageSlack(msg string) []byte {

	if msg == "" {
		msg = slackUserError
	}

	slackTo := struct {
		Text     string `json:"text"`
		Username string `json:"username,omitempty"`
	}{
		msg,
		BotDefaultName,
	}
	var body []byte
	body, err = json.Marshal(slackTo)
	if err != nil {
		body = []byte(strconv.Itoa(http.StatusInternalServerError))
	}

	return body
}

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {
	log.Println("Rimbot invoked.")

	w.Header().Add("Content-Type", "application/json")

	//obtain slacks' webhook

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.Form

	log.Println("text sent to query: ", form.Get("text"))
	//convert webhook values into new outgoing message
	base, target, amount, code := dialogFlow.Query(form.Get("text"))

	log.Println("DialogFlow query output(in rimbot): ", base, "\t", target, "\t", amount, "\t", code)

	if code != http.StatusOK && code != http.StatusPartialContent {
		w.WriteHeader(code)
		return
	} else if code == http.StatusPartialContent { //If Unmarshal fails in Query(). Meaning Clara got confused.
		w.Write(MessageSlack("")) //You fuced up.
	} else { //If everything got parsed correctly.
		errBase := validate.Currency(base)
		errTarget := validate.Currency(target)

		if errBase == nil && errTarget == nil && amount >= 0 { //If valid input for currencyservice.

			currencyTo := map[string]string{
				"baseCurrency":   base,
				"targetCurrency": target,
			}

			body := new(bytes.Buffer)
			err = json.NewEncoder(body).Encode(currencyTo)
			if err != nil {
				w.Write(MessageSlack(""))
				return
			}

			req, err := http.NewRequest(
				http.MethodPost,
				"https://currency-trackr.herokuapp.com/api/latest/", //TODO CHANGE THIS
				ioutil.NopCloser(body),
			)

			log.Printf("Request: %+v", req)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				w.Write(MessageSlack(""))
				return
			}

			// resp, err := http.Post("127.0.0.1:"+os.Getenv("PORT")+"/currency/latest/", "application/json", body)
			// log.Println(body)
			// if err != nil {
			// 	MessageSlack(w, "Post fail")
			// 	return
			// }

			log.Println("respBody: ", resp)
			unParsedRate, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.Write(MessageSlack(""))
				return
			}
			parsedRate, err := strconv.ParseFloat(string(unParsedRate), 64)
			if err != nil {
				w.Write(MessageSlack(""))
				return
			}
			convertedRate := amount * parsedRate
			w.Write(MessageSlack(fmt.Sprintf("%v %v is equal to %v %v. ^^", amount, base, convertedRate, target))) //Temporary outprint

		} else { //If invalid input for currencyservice.
			w.Write(MessageSlack("")) //You fucked up.
			return
		}
	}
}
