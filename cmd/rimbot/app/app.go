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

	switch code {
	case http.StatusOK:
		break
	case http.StatusPartialContent:
		w.Write(MessageSlack("", true)) //You fuced up
		return
	default:
		w.WriteHeader(code)
		return
	}
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
		http.MethodPost, // Posting to the lastest handler of the service.
		os.Getenv("CURRENCY_URI")+"currency/latest/",
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
	w.Write(MessageSlack(fmt.Sprintf("%f %v is equal to %f %v. ^^", amount, base, convertedRate, target), true)) //Temporary outprint

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
	errBase := validate.Currency(base)
	errTarget := validate.Currency(target)
	unparsedAmount := result.Result.Parameters.Amount
	if errBase != nil || errTarget != nil || unparsedAmount == "" {
		code = http.StatusPartialContent
		amount = 1
		return
	}
	amount, err = strconv.ParseFloat(unparsedAmount, 64)
	if err != nil {
		code = http.StatusPartialContent
	}
	return
}
