package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/invoke"

	"github.com/Arxcis/imt2681-assignment2/lib/database"
	"github.com/Arxcis/imt2681-assignment2/lib/httperror"
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/Arxcis/imt2681-assignment2/lib/validate"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// App ...
type App struct {
	Path              string
	Port              string
	CollectionWebhook string
	CollectionFixer   string
	Mongo             database.Mongo
	Currency          []string
}

var err error

// HelloWorld ...
func (app *App) HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

// PostWebhook ...
// POST    /api/v1/subscription/   create a subscription
func (app *App) PostWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Decode webook
	webhook := &types.Webhook{}
	if err = json.NewDecoder(r.Body).Decode(webhook); err != nil {
		httperror.InternalServer(w, err)
		return
	}

	// 2. Validate webhook
	if err = validate.Webhook(webhook, app.Currency); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// 3. Open collection
	cwebhook, err := app.Mongo.OpenC(app.CollectionWebhook)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

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
func (app *App) GetWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := app.Mongo.OpenC(app.CollectionWebhook)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find webhook)
	rawID := mux.Vars(r)["id"]
	if rawID == "" {
		httperror.BadRequest(w, fmt.Errorf("Bad ID"))
		return
	}
	queryID := bson.ObjectIdHex(rawID)
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
func (app *App) GetWebhookAll(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := app.Mongo.OpenC(app.CollectionWebhook)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find all webhooks
	hooks := []types.Webhook{}
	err = cwebhook.Find(nil).All(&hooks)
	//log.Println(hooks)
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
func (app *App) DeleteWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Open collection
	cwebhook, err := app.Mongo.OpenC(app.CollectionWebhook)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

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
func (app *App) GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	// 1. Decode request body
	reqBody := &types.CurrencyIn{}

	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("Malformed body"))
		return
	}

	// 1.5 Validate
	if validate.Currency(reqBody.BaseCurrency, app.Currency) != nil {
		httperror.BadRequest(w, fmt.Errorf("Invalid base currency"))
		return
	}

	if validate.Currency(reqBody.TargetCurrency, app.Currency) != nil {
		httperror.BadRequest(w, fmt.Errorf("Invalid target currency"))
		return
	}

	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		httperror.BadRequest(w, fmt.Errorf("Base currency cannot be equal target currency"))
		return
	}

	// 2. Open collection
	cfixer, err := app.Mongo.OpenC(app.CollectionFixer)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

	// 3. Find latest entry in fixer collection
	fixer := types.FixerIn{}
	err = cfixer.Find(bson.M{"base": reqBody.BaseCurrency}).Sort("-date").One(&fixer)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	// 4. Respond
	fmt.Fprintf(w, "%.2f", fixer.Rates[reqBody.TargetCurrency])
}

// GetAverageCurrency ...
func (app *App) GetAverageCurrency(w http.ResponseWriter, r *http.Request) {

	// 1. Decode request body
	reqBody := &types.CurrencyIn{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("Malformed body"))
		return
	}

	// 1.5 Validate
	if validate.Currency(reqBody.BaseCurrency, app.Currency) != nil {
		httperror.BadRequest(w, fmt.Errorf("Invalid base currency"))
		return
	}

	if validate.Currency(reqBody.TargetCurrency, app.Currency) != nil {
		httperror.BadRequest(w, fmt.Errorf("Invalid target currency"))
		return
	}

	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		httperror.BadRequest(w, fmt.Errorf("Base currency cannot be equal target currency"))
		return
	}

	// 2. Open database
	cfixer, err := app.Mongo.OpenC(app.CollectionFixer)
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find sorted on date descending
	var average float64
	const dayCount int = 3
	fixerArray := []types.FixerIn{}

	err = cfixer.Find(bson.M{"base": reqBody.BaseCurrency}).Sort("-date").Limit(dayCount).All(&fixerArray)
	if err != nil {
		httperror.InternalServer(w, err)
		return
	}

	if len(fixerArray) == 0 {
		httperror.NotFound(w, fmt.Errorf("No fixers found on base currency"))
		return
	}
	// 3. Average last 3 days
	for _, fixer := range fixerArray {
		average += fixer.Rates[reqBody.TargetCurrency]
	}
	average /= float64(dayCount)

	// 4. Respond
	fmt.Fprintf(w, "%.2f", average)
}

// EvaluationTrigger ...
func (app *App) EvaluationTrigger(w http.ResponseWriter, r *http.Request) {

	// 1. Open database
	db, err := app.Mongo.Open()
	if err != nil {
		httperror.ServiceUnavailable(w, err)
		return
	}
	defer app.Mongo.Close()

	// 2. Invoke all webhooks
	client := &http.Client{}
	invoke.Webhooks(client, db.C(app.CollectionWebhook), db.C(app.CollectionFixer))
}
