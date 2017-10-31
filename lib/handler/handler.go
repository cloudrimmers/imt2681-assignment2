package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
	"github.com/Arxcis/imt2681-assignment2/lib/tool"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	webhook.ID = bson.NewObjectId()
	err = db.C("hook").Insert(webhook)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	text, _ := webhook.ID.MarshalText()
	log.Println("ID: ", text)
	w.Write(text)
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	hook := &mytypes.WebhookIn{}

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	queryID := bson.ObjectIdHex(mux.Vars(r)["id"])
	log.Println("GET - hook.id: ", queryID)

	err = db.C("hook").FindId(queryID).One(hook)
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

	hooks := []mytypes.WebhookIn{}
	err = db.C("hook").Find(nil).All(&hooks)
	if err != nil {
		notFound(w, err)
		return
	}

	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
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

	fixer := &mytypes.FixerIn{}
	err = db.C("tick").Find(bson.M{"datestamp": tool.Todaystamp()}).One(fixer)
	if err != nil {
		notFound(w, err)
		return
	}
	fmt.Fprintf(w, "%.2f", fixer.Rates[latest.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
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
		fixer := mytypes.FixerIn{}
		err = db.C("tick").Find(bson.M{"datestamp": tool.Daystamp(i)}).One(&fixer)

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

	hooks := []mytypes.WebhookIn{}
	db.C("hook").Find(nil).All(&hooks)

	for i, hook := range hooks {

		fixer := mytypes.FixerIn{}
		err = db.C("tick").Find(bson.M{"datestamp": tool.Todaystamp()}).One(&fixer)
		if err == nil {
			hook.CurrentRate = fixer.Rates[hook.TargetCurrency]
			go hook.Trigger()
		} else {
			log.Println("webhook ", i, " not triggered..")
		}
	}
}

func serviceUnavailable(w http.ResponseWriter, err error) {
	log.Println("No database connection: ", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println("Collection.Insert() error", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func notFound(w http.ResponseWriter, err error) {
	log.Println("Collection.Find() not found", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println("Http bad request", err.Error())
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
