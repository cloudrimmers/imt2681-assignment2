package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudrimmers/imt2681-assignment3/lib/botNameGenerator"
	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
	"github.com/cloudrimmers/imt2681-assignment3/lib/validate"
)

var err error

const slackUserError = "Sorry, I missed that. Maybe something was vague with what you said? Try using capital letters like this: 'USD', 'GBP'. And numbers like this: '131.5'"
const BotDefaultName = "Rimbot"

func MessageSlack(msg string, generateBotName bool) []byte {

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

	if generateBotName {

		slackTo.Username = botNameGenerator.Generate()
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

	base, target, amount, code := QueryCurrencyConversion(form.Get("text"))

	log.Println("DialogFlow query output(in rimbot): ", base, "\t", target, "\t", amount, "\t", code)

	if code != http.StatusOK && code != http.StatusPartialContent {
		w.WriteHeader(code)
		return
	} else if code == http.StatusPartialContent { //If Unmarshal fails in Query(). Meaning Clara got confused.
		w.Write(MessageSlack("", true)) //You fuced up.
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
				w.Write(MessageSlack("", true))
				return
			}

			req, err := http.NewRequest( //Starts to construct a request.
				http.MethodPost,
				"https://shrouded-journey-80451.herokuapp.com/currency/latest/",
				ioutil.NopCloser(body),
			)

			log.Printf("Request: %+v", req)

			resp, err := http.DefaultClient.Do(req) // Sends request to currencyservice and revieves response.
			if err != nil {
				w.Write(MessageSlack("", true)) //They fucked up.
				return
			}

			log.Println("respBody: ", resp)
			parsedRate, err := ParseFixerResponse(resp.Body)
			if err != nil {
				w.Write(MessageSlack("", true)) //We fucked up.
				return
			}

			convertedRate := amount * parsedRate
			w.Write(MessageSlack(fmt.Sprintf("%v %v is equal to %v %v. ^^", amount, base, convertedRate, target), true)) //Temporary outprint

		} else { //If invalid input for currencyservice.
			w.Write(MessageSlack("", true)) //You fucked up.
			return
		}
	}
}

// ConversionResponse - A representation of the resposnse from DialogFlow
type ConversionResponse struct {
	Result struct {
		//NOTE: If need be, place ADDITIONAL PARAMETERS
		Parameters struct {
			CurrencyOut struct {
				CurrencyName string `json:"currency-name"`
			} `json:"currency-out"`
			CurrencyIn struct {
				CurrencyName string `json:"currency-name"`
			} `json:"currency-in"`
			Amount string `json:"amount"`
		} `json:"parameters"`
	} `json:"result"`
	SessionID string `json:"sessionId"`
}

func (r *ConversionResponse) GetSessionID() string {
	return r.SessionID
}

func QueryCurrencyConversion(text string) (base, target string, amount float64, code int) {
	result := ConversionResponse{}
	code = dialogFlow.Query(text, &result, os.Getenv("ACCESS_TOKEN"))

	base = result.Result.Parameters.CurrencyIn.CurrencyName
	target = result.Result.Parameters.CurrencyOut.CurrencyName
	unparsedAmount := result.Result.Parameters.Amount
	if base == "" || target == "" || unparsedAmount == "" {
		code = http.StatusPartialContent
		amount = 0
		return
	}
	amount, err = strconv.ParseFloat(unparsedAmount, 64)
	if err != nil {
		code = http.StatusPartialContent
	}
	return
}
