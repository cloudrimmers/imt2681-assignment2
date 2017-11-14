package app

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"github.com/gorilla/mux"
)

var APP = App{}

func init() {

	// 2. Initialize the app object
	APP = App{
		FixerioURI:          "https://api.fixer.io/latest",
		CollectionFixerName: "fixer",
		Mongo: database.Mongo{
			Name:    "test",
			URI:     "mongodb://localhost",
			Session: nil,
		},
	}

	// 4. Ensure index to avoid duplicates
	APP.Mongo.EnsureIndex(APP.CollectionFixerName, []string{"date"})

	// 5. Optional seed and log app object
	APP.SeedFixerdata()
	indented, _ := json.MarshalIndent(APP, "", "    ")
	log.Println("App data: ", string(indented))

}

var fixertest = types.FixerIn{
	Base: "EUR",
	Date: "2017-10-22",
	Rates: map[string]float32{
		"NOK": 9.3883,
		"TRY": 4.3751,
		"USD": 1.1761,
		"ZAR": 16.14,
	},
}

func TestFixer2Mongo(t *testing.T) {

	if err = APP.Fixer2Mongo(&fixertest); err != nil {
		t.Log(err.Error())
	}
}

func TestFixerResponse(t *testing.T) {

	testfixerio := func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(fixertest)
		w.Write(bytes)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", testfixerio).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()

	fixer, err := APP.FixerResponse(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Response", fixer)
}
