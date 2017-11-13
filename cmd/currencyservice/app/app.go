package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudrimmers/imt2681-assignment3/lib/database"
)

// App ...
type App struct {
	Port                string
	CollectionFixerName string
	Mongo               database.Mongo
	Currency            []string
}

var err error

// GetLatestCurrency ...
func (app *App) GetLatestCurrency(w http.ResponseWriter, r *http.Request) {

	_ = r.Form["from"]
	_ = r.Form["to"]
	// @TODO - Change to query parameters
	/*
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
		cfixer, err := app.Mongo.OpenC(app.CollectionFixerName)
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
	*/
	fmt.Fprintf(w, "Here is the latest currency")
	log.Println("GET /currency/latest", w.Header().Get("status"))
}
