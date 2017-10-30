package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/Arxcis/imt2681-assignment2/lib/database"

	"github.com/gorilla/mux"

	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
)

// HelloWorld ...
// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

// PostWebhook ...
// POST    /api/v1/subscription/   create a subscription
func PostWebhook(w http.ResponseWriter, r *http.Request) {

	webhook := &mytypes.WebhookIn{}
	_ = json.NewDecoder(r.Body).Decode(webhook)
	fmt.Println(webhook)

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	webhook.ID = bson.NewObjectId()
	db.C("hook").Insert(webhook)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(webhook.ID))
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	w.Header().Add("content-type", "application/json")

	hook := mytypes.WebhookIn{}
	hook.ID = bson.ObjectId(vars["id"])
	data, _ := json.Marshal(hook)
	w.Write(data)

}

// GetWebhookAll ...
func GetWebhookAll(w http.ResponseWriter, r *http.Request) {

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookIn{}

	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	fixer := mytypes.FixerIn{}
	fmt.Fprintf(w, "%.2f", fixer.Rates[latest.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	fmt.Fprintf(w, "%.2f", computeAverage())
}

// EvaluationTrigger ...
func EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	db, err := database.Open()
	if err != nil {

		// Handle error
	}
	defer database.Close()

	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookOut{}
	data, _ := json.Marshal(hooks)
	w.Write(data)
}

func computeAverage() float64 {
	return 8.7
}
