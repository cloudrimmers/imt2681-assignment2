package types

import "gopkg.in/mgo.v2/bson"

// FixerIn ...
/* Example JSON:
{
	"base":"EUR",
    "date":"2017-10-24",
    "rates":{
        "AUD":1.5117,
        "BGN":1.9558,
		.....
		"ZAR":16.14
    }
*/
type FixerIn struct {
	ID        bson.ObjectId      `json:"id" bson:"_id,omitempty"`
	Timestamp string             `bson:"timestamp,omitempty"`
	Base      string             `bson:"base"`
	Date      string             `bson:"date"`
	Rates     map[string]float32 `bson:"rates"`
}

// FixerSeed - a seed data which is meant to seed the database
var FixerSeed = []FixerIn{{
	Base: "EUR",
	Date: "2017-11-01",
	Rates: map[string]float32{
		"AUD": 1.5136,
		"BGN": 1.9558,
		"BRL": 3.815,
		"CAD": 1.4986,
		"CHF": 1.164,
		"CNY": 7.6767,
		"CZK": 25.557,
		"DKK": 7.4415,
		"GBP": 0.87385,
		"HKD": 9.0593,
		"HRK": 7.52,
		"HUF": 311.75,
		"IDR": 15774,
		"ILS": 4.0845,
		"INR": 75.013,
		"JPY": 132.6,
		"KRW": 1291.5,
		"MXN": 22.276,
		"MYR": 4.9136,
		"NOK": 9.461,
		"NZD": 1.6866,
		"PHP": 60.007,
		"PLN": 4.2338,
		"RON": 4.6033,
		"RUB": 67.51,
		"SEK": 9.7535,
		"SGD": 1.5808,
		"THB": 38.505,
		"TRY": 4.4394,
		"USD": 1.1612,
		"ZAR": 16.391,
	},
},
}
