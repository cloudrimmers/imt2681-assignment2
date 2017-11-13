package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"github.com/subosito/gotenv"
	"gopkg.in/mgo.v2/bson"
)

var APP *App
var testid bson.ObjectId

func init() {

	const envpath = "../../../.env"

	// 1. Require .env to be present
	log.Println("Reading .env")
	gotenv.MustLoad(envpath)

	APP = &App{
		Port:                "5555",
		CollectionFixerName: "testfixer",
		Mongo: database.Mongo{
			URI:     os.Getenv("MONGODB_URI"),
			Name:    os.Getenv("MONGODB_NAME"),
			Session: nil,
		},
	}

	// 3. Default values if empty environment
	if APP.Mongo.URI == "" {
		log.Println("No .env present. Using default values")
		APP.Mongo.URI = "mongodb://localhost"
		APP.Mongo.Name = "test"
	}

	log.Println("Seeding DB")
	seedTestDB()

	log.Println("TEST currencyservice initialized...")
}
func TestGetLatestCurrency(t *testing.T) {

	// 1. Setup router
	r := mux.NewRouter()
	r.HandleFunc("/currency/latest", APP.GetLatestCurrency).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + "/currency/latest"
	t.Log("Testing", url)

	// @TODO - change the POST request to a GET request
	// 2. Define table
	table := map[types.CurrencyIn]int{
		//{BaseCurrency: "RAR", TargetCurrency: "NOK"}: http.StatusBadRequest, //@todo - This should be handled as a bad request
		{BaseCurrency: "EUR", TargetCurrency: "EUR"}: http.StatusBadRequest,
		{BaseCurrency: "AAA", TargetCurrency: "EUR"}: http.StatusBadRequest,
		{BaseCurrency: "EUR", TargetCurrency: "AAA"}: http.StatusBadRequest,
		{BaseCurrency: "USD", TargetCurrency: "EUR"}: http.StatusNotFound,
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

func seedTestDB() error {

	collectionFixer, err := APP.Mongo.OpenC(APP.CollectionFixerName)
	if err != nil {
		return err
	}
	defer APP.Mongo.Close()

	collectionFixer.DropCollection()
	collectionFixer.Insert(bson.M{
		"base": "EUR",
		"date": "2017-10-24",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	}, bson.M{
		"base": "EUR",
		"date": "2017-10-23",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	}, bson.M{
		"base": "EUR",
		"date": "2017-10-22",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	})
	return nil
}
