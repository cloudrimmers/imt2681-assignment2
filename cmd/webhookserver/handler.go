package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/httperror"
	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var err error

// HelloWorld ...
// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

// PostWebhook ...
// POST    /api/v1/subscription/   create a subscription
func PostWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Decode webook
	webhook := &types.Webhook{}
	if err = json.NewDecoder(r.Body).Decode(webhook); err != nil {
		httperror.InternalServer(w, err)
		return
	}

	// 2. Validate webhook
	if err = tool.ValidateWebhook(webhook); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// 3. Open collection
	cwebhook, err := database.OpenWebhook()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 4. Insert webhook
	webhook.ID = bson.NewObjectId()
	cwebhook.Insert(webhook)
	if err != nil {
		httperror.InternalServer(w, err)
		return
	}

	// 5. Write response
	w.WriteHeader(http.StatusCreated)
	text, _ := webhook.ID.MarshalText()
	w.Write(text)
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := database.OpenWebhook()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 2. Find webhook
	queryID := bson.ObjectIdHex(mux.Vars(r)["id"])
	hook := &types.Webhook{}

	err = cwebhook.FindId(queryID).One(hook)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	// 3. Marshal and write response
	w.Header().Add("content-type", "application/json")
	data, _ := json.Marshal(hook)
	w.Write(data)

}

// GetWebhookAll ...
func GetWebhookAll(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := database.OpenWebhook()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 2. Find all webhooks
	hooks := []types.Webhook{}
	err = cwebhook.Find(nil).All(&hooks)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	// 3. Marshall and write response
	w.Header().Add("content-type", "application/json")
	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// DeleteWebhook ...
func DeleteWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := database.OpenWebhook()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 2. Log and delete
	queryID := bson.ObjectIdHex(mux.Vars(r)["id"])

	err = cwebhook.RemoveId(queryID)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}
	log.Println("Deleted ", queryID)

	// 3. Implicit OK returned
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	// 1. Decode request body
	reqQuery := &types.CurrencyIn{}
	err := json.NewDecoder(r.Body).Decode(reqQuery)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// 2. Open collection
	cfixer, err := database.OpenFixer()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 3. Find latest entry in fixer collection
	fixer := types.FixerIn{}
	err = cfixer.Find(bson.M{"base": reqQuery.BaseCurrency}).Sort("-date").One(&fixer)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	// 4. Respond
	fmt.Fprintf(w, "%.2f", fixer.Rates[reqQuery.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	// 1. Decode request body
	request := &types.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(request)

	// 2. Open database
	cfixer, err := database.OpenFixer()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 2. Find sorted on date descending
	var average float64
	const dayCount int = 3
	fixerArray := []types.FixerIn{}

	err = cfixer.Find(bson.M{"base": request.BaseCurrency}).Sort("-date").Limit(dayCount).All(&fixerArray)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	// 3. Average last 3 days
	for _, fixer := range fixerArray {
		average += fixer.Rates[request.TargetCurrency]
	}
	average /= float64(dayCount)

	// 4. Respond
	fmt.Fprintf(w, "%.2f", average)
}

// EvaluationTrigger ...
func EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	db, err := database.Open()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer database.Close()

	client := &http.Client{}
	tool.InvokeWebhooks(client, db.C(os.Getenv("COLLECTION_WEBHOOK")), db.C(os.Getenv("COLLECTION_FIXER")))
}
