package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
	"github.com/cloudrimmers/imt2681-assignment3/lib/httperror"
	"github.com/cloudrimmers/imt2681-assignment3/lib/types"
	"github.com/cloudrimmers/imt2681-assignment3/lib/validate"
	"gopkg.in/mgo.v2/bson"
)

// App ...
type App struct {
	Port                string
	CollectionFixerName string
	Mongo               database.Mongo
}

var err error

// GetLatestCurrency ...
func (app *App) GetLatestCurrency(w http.ResponseWriter, r *http.Request) {
	defer log.Println("POST /currency/latest", w.Header().Get("status"))

	var reqBody types.CurrencyIn

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		httperror.BadRequest(w, "json.NewDecoder.Decode()", fmt.Errorf("Malformed body"))
		return
	}

	// 1.5 Validate
	if validate.Currency(reqBody.BaseCurrency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid base currency"))
		return
	}

	if validate.Currency(reqBody.TargetCurrency) != nil {
		httperror.BadRequest(w, "validate.Currency()", fmt.Errorf("Invalid target currency"))
		return
	}

	// 1.6 If two equals currencies, just return 1 == 1 conversion rate
	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		fmt.Fprintf(w, "1.0")
		return
	}

	// 2. Open collection
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		httperror.ServiceUnavailable(w, "app.Mongo.OpenC()", err)
		return
	}
	defer app.Mongo.Close()

	// 3. Find latest entry in fixer collection
	fixer := types.FixerIn{}
	err = collectionFixer.Find(bson.M{"base": reqBody.BaseCurrency}).Sort("-date").One(&fixer)
	if err != nil {
		httperror.NotFound(w, "cfixer.Find()", err)
		return
	}

	// 4. Respond
	fmt.Fprintf(w, "%.2f", fixer.Rates[reqBody.TargetCurrency])

}

// SeedFixerdata ...
// @TODO This is a duplicate of an identical function in fixerworker/ 
//   	  Maybe this should be put in a library - JSolsvik 15.11.17
func (app *App) SeedFixerdata() {

	// 0. Get seed
	seed := types.FixerSeed

	// 1. Open collection
	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer app.Mongo.Close()
	collectionFixer.DropCollection()

	// 2. Insert to database
	// cfixer.DropCollection()
	for _, o := range seed {
		if err = collectionFixer.Insert(o); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}