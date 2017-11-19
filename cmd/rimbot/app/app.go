package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func ParseFixerResponse(body io.ReadCloser) (parsedRate float64, localErr error) {

	unParsedRate, localErr := ioutil.ReadAll(body) // Read all data from request body.
	if localErr != nil {
		return parsedRate, localErr
	}
	parsedRate, localErr = strconv.ParseFloat(string(unParsedRate), 64) // Parse "rate" float from response body.
	if localErr != nil {
		return parsedRate, localErr
	}
	return parsedRate, localErr
}

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {
	log.Println("Rimbot invoked.")

	w.Header().Add("Content-Type", "application/json")

	err = r.ParseForm() //Parse from containing content of message from user.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.Form

	log.Println("text sent to query: ", form.Get("text")) //Print text from form in terminal.
	//convert webhook values into new outgoing message

	/* Query system will soo t\be changed to this convention:
	    func()(base, target, amount, code)...
		  //new response struct
		  // query(&response)
		  result.parameters.cr
	*/
	// base, target, amount, code := func()

	base, target, amount, code := dialogFlow.Query(form.Get("text")) //Gets values from DialogFlow.

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

			currencyTo := map[string]string{ // Request payload to currencyservice.
				"baseCurrency":   base,
				"targetCurrency": target,
			}

			body := new(bytes.Buffer) // Encode request payload to json:
			err = json.NewEncoder(body).Encode(currencyTo)
			if err != nil { // Since values was validated, it "should" be impossible for this to fail.
				w.Write(MessageSlack(""))
				return
			}

			req, err := http.NewRequest( //Starts to construct a request.
				http.MethodPost,
				"https://currency-trackr.herokuapp.com/api/latest/", //TODO CHANGE THIS
				ioutil.NopCloser(body),
			)

			log.Printf("Request: %+v", req)

			resp, err := http.DefaultClient.Do(req) // Sends request to currencyservice and revieves response.
			if err != nil {
				w.Write(MessageSlack("")) //They fucked up.
				return
			}

			log.Println("respBody: ", resp)
			parsedRate, err := ParseFixerResponse(resp.Body)
			if err != nil {
				w.Write(MessageSlack("")) //We fucked up.
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
