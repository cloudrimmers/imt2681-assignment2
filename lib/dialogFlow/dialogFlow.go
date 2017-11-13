package dialogFlow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// query - Queries to DialogFlow use this kind of body
type query struct {
	Query     string   `json:"query"`
	Context   []string `json:"contexts,omitempty"` //may be omitted, context may be implicit
	SessionID int      `json:"sessionId"`
}

// Generate query object to send into dialogFlow
func newQuery(text string) *query {
	return &query{
		Query:     text,
		SessionID: generateSessionID(),
	}
}

func generateSessionID() int {
	return rand.Intn(10000)
}

/*
{
    "id": "a41ed389-151b-415c-ae34-b5093435a5ce",
    "timestamp": "2017-11-13T12:15:01.323Z",
    "lang": "en",
    "result": {
        "source": "agent",
        "resolvedQuery": "Convert 500 NOK to EUR",
        "speech": "",
        "action": "convert",
        "parameters": {
            "currency-out": [
                "EUR"
            ],
            "unit-currency-in": {
                "amount": 500,
                "currency": "NOK"
            }
        },
        "metadata": {
            "inputContexts": [],
            "outputContexts": [],
            "intentName": "convert currency",
            "intentId": "0123783b-f742-4269-9a7b-c9068b66d133",
            "webhookUsed": "false",
            "webhookForSlotFillingUsed": "false",
            "contexts": []
        },
        "score": 1
    },
    "status": {
        "code": 200,
        "errorType": "success"
    },
    "sessionId": "Bois"
}
*/

// Response - A representation of the response from DialogFlow
type Response struct {
	Query  string `json:"query"`
	Result struct {
		//ADDITIONAL PARAMETERS
		Parameters struct { //These may vary... should we use a map perhaps
			CurrencyOut struct {
				CurrencyName string `json:"currency-name,omitempty"`
			} `json:"currency-out"`
			CurrencyIn struct {
				CurrencyName string `json:"currency-name"`
			} `json:"currency-in"`
			Amount float64 `json:"amount"`
		} `json:"parameters"`
	} `json:"result"`
	Status struct {
		Code  int    `json:"code"`
		Error string `json:"errorType"`
	} `json:"status"`
	SessionID string `json:"sessionId"`
}

const dialogFlowRoot = "https://api.dialogflow.com/v1/query?v=" ////NOTE: protocol number is "required", consider adding it

// Protocols: https://dialogflow.com/docs/reference/agent/#protocol_version
const (
	ProtocolBase    = 20150910
	ProtocolNumeric = 20170712
)

//Query DialogFlow for a conversion
func Query(queryText string) (responseObject Response, statusCode int) {
	responseObject = Response{} //prepare responseObject

	query, err := json.Marshal(newQuery(queryText))
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}
	fmt.Printf("%+v\n", query) //Print the body that will be sent to DialigFlow.

	//Construct a request with our query object
	req, err := http.NewRequest(
		http.MethodPost,
		dialogFlowRoot+strconv.Itoa(ProtocolNumeric),
		ioutil.NopCloser(bytes.NewBuffer(query)),
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	//Add authorization token to head. Identifies agent in dialogflow.
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req) //Execute request.
	if err != nil {
		statusCode = http.StatusFailedDependency //NOTE: is this right? - yes it is!
		return
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	err = json.Unmarshal(respBody, &responseObject)
	if err != nil {
		log.Println(err)
		log.Printf("failed unmarshalling response:\n%+v", responseObject)
		statusCode = http.StatusInternalServerError
		return
	}
	statusCode = responseObject.Status.Code
	return
}
