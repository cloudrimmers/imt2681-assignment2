package invoke

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"gopkg.in/mgo.v2/bson"
)

func TestWebhooks(t *testing.T) {

	mongo := database.Mongo{Name: "test", URI: "127.0.0.1:33017", Session: nil}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "!!! Invoked! !!!")
	}))
	defer ts.Close()
	db, err := mongo.Open()
	if err != nil {
		t.Error(err.Error())
	}

	cWebhook := db.C("testhook")
	cFixer := db.C("testfixer")
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
	})

	cWebhook.Insert(bson.M{
		"webhookURL":      ts.URL,
		"baseCurrency":    "EUR",
		"targetCurrency":  "NOK",
		"minTriggerValue": 9.0,
		"maxTriggerValue": 9.9,
	}, bson.M{
		"webhookURL":      ts.URL,
		"baseCurrency":    "EUR",
		"targetCurrency":  "NOK",
		"minTriggerValue": 9.0,
		"maxTriggerValue": 9.9,
	})

	client := http.Client{}
	Webhooks(&client, cWebhook, cFixer)

}
