package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"gopkg.in/mgo.v2/bson"
)

var APP *App
var testid bson.ObjectId

func init() {

	APP = &App{
		Port:                "5555",
		CollectionFixerName: "testfixer",
		Mongo: database.Mongo{
			URI:     "mongodb://localhost",
			Name:    "test",
			Session: nil,
		},
	}
	APP.SeedFixerdata()
	log.Println("TEST currencyservice initialized...")
}

func TestGetLatestCurrency(t *testing.T) {

	// 1. Setup router
	r := mux.NewRouter()
	r.HandleFunc("/currency/latest/", APP.GetLatestCurrency).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + "/currency/latest/"
	t.Log("Testing", url)

	// @TODO - change the POST request to a GET request
	// 2. Define table
	table := map[types.CurrencyIn]int{
		//{BaseCurrency: "RAR", TargetCurrency: "NOK"}: http.StatusBadRequest, //@todo - This should be handled as a bad request
		{BaseCurrency: "EUR", TargetCurrency: "EUR"}: http.StatusOK,
		{BaseCurrency: "AAA", TargetCurrency: "EUR"}: http.StatusBadRequest,
		{BaseCurrency: "EUR", TargetCurrency: "AAA"}: http.StatusBadRequest,
		{BaseCurrency: "USD", TargetCurrency: "EUR"}: http.StatusOK,
		{BaseCurrency: "EUR", TargetCurrency: "NOK"}: http.StatusOK,
	}

	// 3. Run tests
	for postBody, wantedStatus := range table {
		byteBody, _ := json.Marshal(postBody)
		t.Run(string(byteBody), func(t *testing.T) {

			resp, err := http.Post(url, "application/json", bytes.NewReader(byteBody))
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != wantedStatus {
				t.Fatalf("Wrong status code. Got %d want %d", resp.StatusCode, wantedStatus)
			}
		})
	}

	malformedString := "sdjføalsjfløsajdfløjslødj"
	t.Run("Malformed request: "+malformedString, func(t *testing.T) {
		resp, _ := http.Post(url, "application/json", ioutil.NopCloser(bytes.NewBufferString(malformedString)))
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Wrong status code. Got %d want %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

}
