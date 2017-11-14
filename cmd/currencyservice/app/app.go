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

	if reqBody.BaseCurrency == reqBody.TargetCurrency {
		httperror.BadRequest(w, "basecurrency==targetcurrency", fmt.Errorf("Base currency cannot be equal target currency"))
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

// SeedTestDB ...
func (app *App) SeedTestDB() error {

	collectionFixer, err := app.Mongo.OpenC(app.CollectionFixerName)
	if err != nil {
		return err
	}
	defer app.Mongo.Close()

	collectionFixer.DropCollection()
	collectionFixer.Insert(bson.M{
		"base": "EUR",
		"date": "2017-10-24",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	}, bson.M{
		"base": "EUR",
		"date": "2017-10-23",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	}, bson.M{
		"base": "EUR",
		"date": "2017-10-22",
		"rates": map[string]float64{
			"NOK": 9.3883,
			"TRY": 4.3751,
			"USD": 1.1761,
			"ZAR": 16.14,
		},
	})
	return nil
}
