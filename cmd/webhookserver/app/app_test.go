package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/subosito/gotenv"
	"gopkg.in/mgo.v2/bson"
)

var APP *App
var testid bson.ObjectId

func init() {
	const envpath = "../../../.env"
	const configpath = "../../../config/currency.json"

	log.Println("Reading ", envpath)
	gotenv.MustLoad(envpath)
	log.Println("Done with ", envpath)

	APP = &App{
		Path:              "/api/test",
		Port:              "5555",
		CollectionWebhook: "testhook",
		CollectionFixer:   "testfixer",
		Mongo: database.Mongo{
			Name:    os.Getenv("MONGODB_NAME"),
			URI:     os.Getenv("MONGODB_URI"),
			Session: nil,
		},
		Currency: func() []string {
			log.Println("Reading " + configpath)
			data, err := ioutil.ReadFile(configpath)
			if err != nil {
				panic(err.Error())
			}
			var currency []string
			if err = json.Unmarshal(data, &currency); err != nil {
				panic(err.Error())
			}
			log.Println("Done with " + configpath)
			return currency
		}(),
	}
	// @verbose
	// indented, _ := json.MarshalIndent(APP, "", "    ")
	// log.Println(string(indented))
	log.Println("Webhookserver initialized...")

	log.Println("Reseeding DB")
	reseedDB()
	log.Println("Done with DB")
}

func TestHelloWorld(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(APP.HelloWorld))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err.Error())
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Printf("%s", greeting)
}

func TestPostWebhook(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(APP.PostWebhook))
	defer ts.Close()
	json, err := json.Marshal(types.Webhook{
		WebhookURL:      ts.URL,
		BaseCurrency:    "EUR",
		TargetCurrency:  "NOK",
		MinTriggerValue: 7.7,
		MaxTriggerValue: 9.9,
	})

	res, err := http.Post(ts.URL, "application/json", bytes.NewReader(json))
	if err != nil {
		t.Error(err.Error())
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Printf("%s\n", greeting)
}
func TestGetWebhook(t *testing.T) {

	r := mux.NewRouter()
	r.HandleFunc(APP.Path+"/webhook/{id}", APP.GetWebhook).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()

	table := map[string]int{
		"dfdfdfd":                  http.StatusBadRequest,
		"45cbc4a0e4123f6920000002": http.StatusNotFound,
		testid.Hex():               http.StatusOK,
	}

	for id, status := range table {
		t.Run(id, func(t *testing.T) {
			url := ts.URL + APP.Path + "/webhook/" + id
			resp, err := http.Get(url)
			if err != nil {
				t.Fatal(err)
			}
			if s := resp.StatusCode; s != status {
				t.Fatalf("Wrong status code. Got %d want %d", s, status)
			}
		})
	}
}
func TestGetWebhookAll(t *testing.T) {

	r := mux.NewRouter()
	r.HandleFunc(APP.Path+"/webhook/", APP.GetWebhookAll).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + APP.Path + "/webhook/"
	t.Log("Testing", url)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code. Got %d want %d", resp.StatusCode, http.StatusOK)
	}
}
func TestGetLatestCurrency(t *testing.T) {

	// 1. Setup router
	r := mux.NewRouter()
	r.HandleFunc(APP.Path+"/currency/latest", APP.GetLatestCurrency).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + APP.Path + "/currency/latest"
	t.Log("Testing", url)

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
func TestGetAverageCurrency(t *testing.T) {
	// 1. Setup router
	r := mux.NewRouter()
	r.HandleFunc(APP.Path+"/currency/average", APP.GetAverageCurrency).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := ts.URL + APP.Path + "/currency/average"
	t.Log("Testing", url)
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
}
func TestEvaluationTrigger(t *testing.T) {}

func reseedDB() error {

	db, err := APP.Mongo.Open()
	if err != nil {
		return err
	}
	defer APP.Mongo.Close()

	cWebhook := db.C(APP.CollectionWebhook)
	cFixer := db.C(APP.CollectionFixer)

	cFixer.DropCollection()
	cWebhook.DropCollection()

	cFixer.Insert(bson.M{
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
	testid = bson.NewObjectId()
	cWebhook.Insert(bson.M{
		"_id":             testid,
		"webhookURL":      "127.0.0.1:5555",
		"baseCurrency":    "EUR",
		"targetCurrency":  "NOK",
		"minTriggerValue": 9.0,
		"maxTriggerValue": 9.9,
	}, bson.M{
		"webhookURL":      "127.0.0.1:5555",
		"baseCurrency":    "EUR",
		"targetCurrency":  "NOK",
		"minTriggerValue": 9.0,
		"maxTriggerValue": 9.9,
	})
	return nil
}
