package payload

// Ticker ...
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
}
*/
type Ticker struct {
	Base  string
	Date  string
	Rates map[string]float64
}
