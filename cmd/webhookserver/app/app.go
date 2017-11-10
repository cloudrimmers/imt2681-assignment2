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
// POST    apipath/webhook/   create a subscription
func (app *App) PostWebhook(w http.ResponseWriter, r *http.Request) {

	// 1. Decode webook
	webhook := types.Webhook{}

	err = json.NewDecoder(r.Body).Decode(&webhook)

	
	if err != nil {
		httperror.InternalServer(w, "json.NewDecoder.Decode() ", err)
		return
	}

	log.Println("!!!! NEW WEBHOOK ", webhook)

	// 2. Validate webhook
	if err = validate.NewWebhook(&webhook, app.Currency); err != nil {
		httperror.BadRequest(w, "validate.Webhook()", err)
		return
	}

	// 3. Open collection
	cwebhook, err := app.Mongo.OpenC(app.CollectionWebhook)
	if err != nil {
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 4. Insert webhook
	webhook.ID = bson.NewObjectId()
	err = cwebhook.Insert(webhook)
	if err != nil {
		httperror.InternalServer(w, "cwebhook.Insert()", err)
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
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find webhook)
	rawID := mux.Vars(r)["id"]

	log.Println("!!!! GETTTING ID  ", rawID)

	// @HACK - MAD HACK JUST TO COMPLY WITH THE ASSIGNMENT 2 routing shecma
	//   - Redirecting to the /evaluationtrigger - handler
	if rawID == "evaluationtrigger" {
		App.EvaluationTrigger(w, r)
		return
	}
	// @HACK END

	if !bson.IsObjectIdHex(rawID) || rawID == "" {
		httperror.BadRequest(w, "bson.IsObjectIdHex()", fmt.Errorf("Bad objectID %s", rawID))
		return
	}
	queryID := bson.ObjectIdHex(rawID)
	hook := &types.Webhook{}

	err = cwebhook.FindId(queryID).One(hook)
	if err != nil {
		httperror.NotFound(w, "cwebhook.FindId()", err)
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
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find all webhooks
	hooks := []types.Webhook{}
	err = cwebhook.Find(nil).All(&hooks)
	if err != nil {
		httperror.NotFound(w, "cwebhook.Find()", err)
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
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 2. Log and delete
	queryID := bson.ObjectIdHex(mux.Vars(r)["id"])

	err = cwebhook.RemoveId(queryID)
	if err != nil {
		httperror.NotFound(w, "cwebhook.RemoveID()", err)
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
		httperror.BadRequest(w, "json.NewDecoder.Decode()", fmt.Errorf("Malformed body"))
		return
	}

	// 1.5 Validate
	if validate.Currency(reqBody.BaseCurrency, app.Currency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid base currency"))
		return
	}

	if validate.Currency(reqBody.TargetCurrency, app.Currency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid target currency"))
		return
	}

	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		httperror.BadRequest(w, "basecurrency==targetcurrency", fmt.Errorf("Base currency cannot be equal target currency"))
		return
	}

	// 2. Open collection
	cfixer, err := app.Mongo.OpenC(app.CollectionFixer)
	if err != nil {
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 3. Find latest entry in fixer collection
	fixer := types.FixerIn{}
	err = cfixer.Find(bson.M{"base": reqBody.BaseCurrency}).Sort("-date").One(&fixer)
	if err != nil {
		httperror.NotFound(w, "cfixer.Find()", err)
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
		httperror.BadRequest(w, "json.NewDecoder.Decode()", fmt.Errorf("Malformed body"))
		return
	}

	// 1.5 Validate
	if validate.Currency(reqBody.BaseCurrency, app.Currency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid base currency"))
		return
	}

	if validate.Currency(reqBody.TargetCurrency, app.Currency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid target currency"))
		return
	}

	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		httperror.BadRequest(w, "basecurrency==targetcurrency", fmt.Errorf("Base currency cannot be equal target currency"))
		return
	}

	// 2. Open database
	cfixer, err := app.Mongo.OpenC(app.CollectionFixer)
	if err != nil {
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 2. Find sorted on date descending
	var average float64
	const dayCount int = 3
	fixerArray := []types.FixerIn{}

	err = cfixer.Find(bson.M{"base": reqBody.BaseCurrency}).Sort("-date").Limit(dayCount).All(&fixerArray)
	if err != nil {
		httperror.InternalServer(w, "cfixer.Find()", err)
		return
	}

	if len(fixerArray) == 0 {
		httperror.NotFound(w, "len(fixerArray) == 0", fmt.Errorf("No fixers found on base currency"))
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
		httperror.ServiceUnavailable(w, "app.Mongo.Open()", err)
		return
	}
	defer app.Mongo.Close()

	// 2. Invoke all webhooks
	client := &http.Client{}
	invoke.Webhooks(client, db.C(app.CollectionWebhook), db.C(app.CollectionFixer))
}
