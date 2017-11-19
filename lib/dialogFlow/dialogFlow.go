package dialogFlow

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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

//QueryOut - Generalized result of a query
type QueryOut interface {
	GetSessionID() string
}

type status struct {
	Code         int    `json:"code"`
	Error        string `json:"errorType"`
	ErrorDetails string `json:"errorDetails, omitempty"`
}
type queryStatus struct {
	*status `json:"status"`
}

const dialogFlowRoot = "https://api.dialogflow.com/v1/query?v=" //NOTE: protocol number is "required", consider adding it

// Protocols: https://dialogflow.com/docs/reference/agent/#protocol_version
const (
	ProtocolBase    = 20150910
	ProtocolNumeric = 20170712 //NOTE we'll not take use of the numeric protocol as it's inconsistent when taking use of DialogFlows default values
	// Specifically 'amount' of type @sys.number is not sent as a number in the json-body when set through default values, but rather a string.
)

type requester func(req *http.Request) (resp *http.Response, err error)

func doQuery(qry *query, rq requester, result QueryOut, token string) (statusCode int) {
	encodedQuery, err := json.Marshal(qry)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	//Construct a request with our query object
	req, err := http.NewRequest(
		http.MethodPost,
		dialogFlowRoot+strconv.Itoa(ProtocolBase),
		ioutil.NopCloser(bytes.NewBuffer(encodedQuery)),
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	//Add authorization token to head. Identifies agent in dialogflow.
	req.Header.Add("Authorization", "Bearer "+token)
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

	qryStatus := queryStatus{}

	err = json.Unmarshal(respBody, &qryStatus)
	if err != nil {
		log.Printf("Error:\n%v\nFailed unmarshalling response:\n%+v", err, qryStatus)
	}

	statusCode = qryStatus.Code

	if qryStatus.Code != http.StatusOK {
		log.Printf("Dialogflow Error: %+v", qryStatus.status)
		return
	}

	err = json.Unmarshal(respBody, &result) // NOTE: err might be ignored at this point
	if err != nil {
		log.Println(err)
		statusCode = http.StatusPartialContent
		return
	}
	log.Printf("%+v", result)

	if result.GetSessionID() != qry.SessionID {
		statusCode = http.StatusUnauthorized
		return
	}
	return
}

/*
Query - Queries a DialogFlow agent with the given token
Statuses:
InternalServerError
FailedDependency
PartialContent
OK
Unauthorized
*/
func Query(queryText string, result QueryOut, token string) (statusCode int) {
	qry := newQuery(queryText)
	return doQuery(qry, http.DefaultClient.Do, result, token)
}
