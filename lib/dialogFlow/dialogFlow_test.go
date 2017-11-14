package dialogFlow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

const dialogFlowSampleFile = "./testData/dialogFlowSample.json"

func DoRequestHere(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Header = http.Header{}
	resp.Header.Add("Content-Type", "application/json")
	data, err := ioutil.ReadFile(dialogFlowSampleFile)
	if err != nil {
		panic(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
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

func TestQuery(t *testing.T) {

	parser := func(in, out string, amount float64) (qry string, rq requester) {
		qry = fmt.Sprintf("%v %v to %v", amount, in, out)
		rq = func(req *http.Request) (resp *http.Response, err error) {
			resp = &http.Response{}
			resp.StatusCode = http.StatusOK
			resp.Header = http.Header{}
			resp.Header.Add("Content-Type", "application/json")
			data, err := ioutil.ReadFile(dialogFlowSampleFile)
			if err != nil {
				panic(err)
			}
			responseData := Response{}
			json.Unmarshal(data, &responseData)
			responseData.Query = fmt.Sprintf("%v %v to %v", amount, in, out)
			responseData.Result.Parameters.CurrencyIn.CurrencyName = in
			responseData.Result.Parameters.CurrencyOut.CurrencyName = out
			responseData.Result.Parameters.Amount = amount
			raw, err := json.Marshal(responseData)
			if err != nil {
				panic(err)
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
			return
		}
		return
	}
	// type generator
	// tests := []struct {
	// 	name string
	// 	gen  struct {
	// 		qry string
	// 		rq  requester
	// 	}
	// }{
	// 	// {"Test 1", parser("Nok", "GBP", 23.56)},
	// 	// {"Test 2", parser("Nok", "GBP", 23.56)},
	// 	// {"Test 3", parser("Nok", "GBP", 23.56)},
	// }
	//
	// for test := range tests {
	//
	// }
	doQuery(parser("NOK", "USD", 34.543))
}
