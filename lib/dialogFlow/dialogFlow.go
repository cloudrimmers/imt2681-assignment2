package dialogFlow

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// query - Queries to DialogFlow use this kind of body:
type query struct { //The required fields are listed in the DialogFlow docs here: https://dialogflow.com/docs/reference/agent/query
	Language  string   `json:"lang"` // If they change the required fields, it might cause DialogFlow to respond with BadRequest.
	Query     string   `json:"query"`
	Contexts  []string `json:"contexts,omitempty"` //may be omitted, context may be implicit
	SessionID string   `json:"sessionId"`
}

// Generate query object to send into dialogFlow
func newQuery(text string, contexts ...string) *query {
	qry := new(query)
	qry.Language = "en" //Required by Dialog flow since 14-11-17 :/
	qry.Query = text
	qry.Contexts = contexts
	qry.SessionID = generateSessionID()
	return qry
}

func generateSessionID() string {
	return strconv.Itoa(rand.Intn(10000))
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
		//NOTE: If need be, place ADDITIONAL PARAMETERS
		Parameters struct {
			CurrencyOut struct {
				CurrencyName string `json:"currency-name,omitempty"`
			} `json:"currency-out"`
			CurrencyIn struct {
				CurrencyName string `json:"currency-name,omitempty"`
			} `json:"currency-in"`
			Amount float64 `json:"amount,omitempty"`
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

type requester func(req *http.Request) (resp *http.Response, err error)

func doQuery(queryText string, rq requester) (base string, target string, amount float64, statusCode int) {
	responseObject := Response{} //prepare responseObject

	query := newQuery(queryText)

	encodedQuery, err := json.Marshal(query)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	//Construct a request with our query object
	req, err := http.NewRequest(
		http.MethodPost,
		dialogFlowRoot+strconv.Itoa(ProtocolNumeric),
		ioutil.NopCloser(bytes.NewBuffer(encodedQuery)),
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	//Add authorization token to head. Identifies agent in dialogflow.
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	log.Printf("%+v", req)
	response, err := rq(req) //Execute request.
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
		statusCode = http.StatusPartialContent
		return
	}
	log.Printf("%+v", responseObject)
	// DANGER!!! - someone
	if responseObject.SessionID != query.SessionID {
		statusCode = http.StatusUnauthorized
		responseObject = Response{}
		return
	}
	base = responseObject.Result.Parameters.CurrencyIn.CurrencyName
	target = responseObject.Result.Parameters.CurrencyOut.CurrencyName
	amount = responseObject.Result.Parameters.Amount

	statusCode = responseObject.Status.Code
	return
}

//Query DialogFlow for a conversion
func Query(queryText string) (base string, target string, amount float64, statusCode int) {
	return doQuery(queryText, http.DefaultClient.Do)
}
