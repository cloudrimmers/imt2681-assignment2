package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

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

	database, _ := database.Get()
	webhook.ID = bson.NewObjectId()
	database.C("hook").Insert(webhook)

	id := strconv.Itoa(rand.Int())
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// ERROR HANDLING

	w.Header().Add("content-type", "application/json")

	hook := mytypes.WebhookIn{}
	hook.ID = vars["id"]
	data, _ := json.Marshal(hook)
	w.Write(data)

}

// GetWebhookAll ...
func GetWebhookAll(w http.ResponseWriter, r *http.Request) {

	// ERROR HANDLING

	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookIn{}

	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	// ERROR HANDLING

	fixer := mytypes.FixerIn{}
	fmt.Fprintf(w, "%.2f", fixer.Rates[latest.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	// ERROR HANDLING

	fmt.Fprintf(w, "%.2f", computeAverage())
}

// EvaluationTrigger ...
func EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	// ERROR HANDLING

	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookOut{}
	data, _ := json.Marshal(hooks)
	w.Write(data)
}

func computeAverage() float64 {
	return 8.7
}
