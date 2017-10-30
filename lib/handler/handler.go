package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
	"github.com/gorilla/mux"
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
	fmt.Println(webhook)

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
	w.Write([]byte(webhook.ID))
}

// GetWebhook ...
func GetWebhook(w http.ResponseWriter, r *http.Request) {

	hook := mytypes.WebhookIn{}
	hook.ID = bson.ObjectId(mux.Vars(r)["id"])

	db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()

	query := db.C("hook").FindId(hook.ID)
	err = query.One(hook)
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

	/*	db, err := database.Open()
		if err != nil {
			serviceUnavailable(w, err)
			return
		}
		defer database.Close()
	*/
	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookIn{}

	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// GetLatestCurrency ...
func GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	/*	db, err := database.Open()
		if err != nil {
			serviceUnavailable(w, err)
			return
		}
		defer database.Close()
	*/
	fixer := mytypes.FixerIn{}
	fmt.Fprintf(w, "%.2f", fixer.Rates[latest.TargetCurrency])
}

// GetAverageCurrency ...
func GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	latest := &mytypes.CurrencyIn{}
	_ = json.NewDecoder(r.Body).Decode(latest)

	/*db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()
	*/
	fmt.Fprintf(w, "%.2f", computeAverage())
}

// EvaluationTrigger ...
func EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	/*db, err := database.Open()
	if err != nil {
		serviceUnavailable(w, err)
		return
	}
	defer database.Close()
	*/
	w.Header().Add("content-type", "application/json")

	hooks := []mytypes.WebhookOut{}
	data, _ := json.Marshal(hooks)
	w.Write(data)
}

// Private
func computeAverage() float64 {
	return 8.7
}

func serviceUnavailable(w http.ResponseWriter, err error) {
	log.Println("No database connection: ", err.Error())
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println("Collection.Insert error", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func notFound(w http.ResponseWriter, err error) {
	log.Println("Hook not found", err.Error())
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
