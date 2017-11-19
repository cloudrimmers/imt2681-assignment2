package dialogFlow

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/cloudrimmers/imt2681-assignment3/lib/reflectUtil"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"github.com/subosito/gotenv"
)

const dialogFlowTestSamples = "./testData/"

func DoRequestHere(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Header = http.Header{}
	resp.Header.Add("Content-Type", "application/json")
	data, err := ioutil.ReadFile(dialogFlowTestSamples + "OK.json")
	if err != nil {
		panic(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return
}

func DoUnauthorizedRequestHere(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}
	resp.StatusCode = http.StatusUnauthorized
	resp.Header = http.Header{}
	resp.Header.Add("Content-Type", "application/json")
	data, err := ioutil.ReadFile(dialogFlowTestSamples + "Unauthorized.json")
	if err != nil {
		panic(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return
}

func DoBadRequestHere(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}
	resp.StatusCode = http.StatusBadRequest
	resp.Header = http.Header{}
	resp.Header.Add("Content-Type", "application/json")
	data, err := ioutil.ReadFile(dialogFlowTestSamples + "BadRequest.json")
	if err != nil {
		panic(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	return
}

func DoRequestError(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}
	err = fmt.Errorf("Something went wrong u big dumdum")
	return
}

func TestNewQuery(t *testing.T) {
	//SETUP
	tests := []string{
		"Convert 200 bucks to nok",
		"Convert 200 bucks to nok",
		"¥ of 500$",
		"100 Rupee in Danish krone",
	}
	for i, test := range tests {

		qry := newQuery(test)
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test != qry.Query {
				wanted := query{
					Query:     test,
					Contexts:  qry.Contexts,
					SessionID: qry.SessionID,
				}
				t.Errorf("newQuery() = %v want %v", qry, wanted)
			}
			if len(qry.SessionID) <= 0 {
				wanted := query{
					Query:     test,
					Contexts:  qry.Contexts,
					SessionID: "NOT EMPTY",
				}
				t.Errorf("newQuery() = %v want %v", qry, wanted)
			}
		})
	}
	//Teardown
}

type queryRequesterGen struct {
	qry string
	rq  requester
}

func TestDialogFlowDependencyFail(t *testing.T) {

	testOut := types.Response{}
	status := doQuery(nil, DoRequestError, &testOut, "")

	if status != http.StatusFailedDependency {
		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerNameInTest(), status, http.StatusFailedDependency)
	}

	t.Logf("%+v", testOut)
}

func TestBadRequest(t *testing.T) {

	testOut := types.Response{}
	status := doQuery(nil, DoBadRequestHere, &testOut, "")

	if status != http.StatusBadRequest {
		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerNameInTest(), status, http.StatusBadRequest)
	}

	t.Logf("%+v", testOut)
}

func TestUnauthorizedRequest(t *testing.T) {
	testOut := types.Response{}
	status := doQuery(nil, DoUnauthorizedRequestHere, &testOut, "")
	if status != http.StatusUnauthorized {
		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerNameInTest(), status, http.StatusUnauthorized)
	}
	t.Logf("%+v", testOut)
}

func TestValidQuery(t *testing.T) {
	testOut := types.Response{}
	qry := &query{
		Language:  "en",
		Query:     "Convert 200 nok to USD",
		SessionID: "1",
		Contexts:  nil,
	}
	status := doQuery(qry, DoRequestHere, &testOut, "")

	if status != http.StatusOK {
		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerName(), status, http.StatusOK)
	}
	t.Logf("%+v", testOut)
}

func TestWrongSessionID(t *testing.T) {
	testOut := types.Response{}
	qry := &query{
		Language:  "en",
		Query:     "Convert 200 nok to USD",
		SessionID: "NOT THE SESSION ID YOU'RE LOOKING FOR",
		Contexts:  nil,
	}
	status := doQuery(qry, DoRequestHere, &testOut, "")

	if status != http.StatusUnauthorized {
		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerName(), status, http.StatusUnauthorized)
	}

	t.Logf("%+v", testOut)
}

//
// func TestPartialContent(t *testing.T) {
// 	testOut := struct{}{}
//
//
// 	status := doQuery(nil, DoRequestHere, &testOut, "")
//
// 	if status != http.StatusPartialContent {
// 		t.Errorf("%s = %v, want %v", reflectUtil.GetCallerName(), status, http.StatusPartialContent)
// 	}
//
// 	t.Logf("%+v", testOut)
// }

//NOTE: Not really neccessary, but just to confirm nothing is wrong with our mock of dialogflow
func TestReal(t *testing.T) {
	gotenv.MustLoad("../../cmd/rimbot/.env")
	testOut := types.Response{}
	status := Query("Nok to Eur", &testOut, os.Getenv("ACCESS_TOKEN"))
	if status != http.StatusOK {
		t.Errorf("Actually connecting to dialogflow fails, Check if API is changed")
		t.Errorf("%+v", testOut)
	}
	t.Logf("%+v", testOut)
}
