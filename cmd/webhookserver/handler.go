package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *types.WebConfig = (&types.WebConfig{}).Load()

// HelloWorld ...
// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

// PostWebhook ...
// POST    /api/v1/subscription/   create a subscription
func PostWebhook(w http.ResponseWriter, r *http.Request) {

	// 0. Decode webook
	webhook := &types.Webhook{}
	if errDecode := json.NewDecoder(r.Body).Decode(webhook); errDecode != nil {
		internalServerError(w, errDecode)
		return
	}

	// 1. Open database
	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	// 2. Validate webhook
	if err = tool.ValidateWebhook(webhook, config); err != nil {
		badRequest(w, err)
		return
	}

	// 3. Insert webhook
	webhook.ID = bson.NewObjectId()
	err = db.C(config.CollectionWebhook).Insert(webhook)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// 4. Write header
	w.WriteHeader(http.StatusCreated)
	text, _ := webhook.ID.MarshalText()
	w.Write(text)
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	hook := &types.Webhook{}

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	queryID := bson.ObjectIdHex(mux.Vars(r)["id"])
	log.Println("GET - hook.id: ", queryID)

	err = db.C(config.CollectionWebhook).FindId(queryID).One(hook)
	if err != nil {
		notFound(w, err)
		return
	}

	w.Header().Add("content-type", "application/json")

	data, _ := json.Marshal(hook)
	w.Write(data)

}

// GetWebhookAll ...
func GetWebhookAll(w http.ResponseWriter, r *http.Request) {

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	w.Header().Add("content-type", "application/json")

	hooks := []types.Webhook{}
	err = db.C(config.CollectionWebhook).Find(nil).All(&hooks)
	if err != nil {
		notFound(w, err)
		return
	}

	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &types.CurrencyIn{}
	err := json.NewDecoder(r.Body).Decode(latest)
	if err != nil {
		badRequest(w, err)
		return
	}
	db := &mgo.Database{}
	db, err = database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	fixer := &types.FixerIn{}
	err = db.C(config.CollectionTick).Find(bson.M{"datestamp": tool.Todaystamp()}).One(fixer)
	if err != nil {
		notFound(w, err)
		return
	}
	fmt.Fprintf(w, "%.2f", fixer.Rates[latest.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &types.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	// @note look 3 days back
	const dayCount int = 3
	var average float64

	for i := 0; i < dayCount; i++ {
		fixer := types.FixerIn{}
		err = db.C(config.CollectionTick).Find(bson.M{"datestamp": tool.Daystamp(i)}).One(&fixer)

		log.Println("i: ", i, "data: ", tool.Daystamp(i), fixer, latest.BaseCurrency)

		if err != nil {
			notFound(w, err)
			return
		}
		average += fixer.Rates[latest.TargetCurrency]
	}
	average /= float64(dayCount)

	fmt.Fprintf(w, "%.2f", average)
}

// EvaluationTrigger ...
// Invoke all wehooks
func EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	hooks := []types.Webhook{}
	db.C(config.CollectionWebhook).Find(nil).All(&hooks)

	for i, hook := range hooks {

		fixer := types.FixerIn{}
		err = db.C(config.CollectionTick).Find(bson.M{"datestamp": tool.Todaystamp()}).One(&fixer)
		if err == nil {
			hook.CurrentRate = fixer.Rates[hook.TargetCurrency]
			go hook.Trigger()
		} else {
			log.Println("webhook ", i, " not triggered..")
		}
	}
}
