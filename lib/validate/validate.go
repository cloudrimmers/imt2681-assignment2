package validate

import (
	"fmt"
	"regexp"
)

var dbregex = "(mongodb:\\/\\/localhost(:\\d{1,5})?)|(127.0.0.1(:\\d{1,5})?)|(mongodb:\\/\\/heroku_(.*):(.*)@(.*).mlab.com:(\\d){1,5}\\/(.*))"

var currencies = [...]string{
	"EUR",
	"AUD",
	"BGN",
	"BRL",
	"CAD",
	"CHF",
	"CNY",
	"CZK",
	"DKK",
	"GBP",
	"HKD",
	"HRK",
	"HUF",
	"IDR",
	"ILS",
	"INR",
	"JPY",
	"KRW",
	"MXN",
	"MYR",
	"NOK",
	"NZD",
	"PHP",
	"PLN",
	"RON",
	"RUB",
	"SEK",
	"SGD",
	"THB",
	"TRY",
	"USD",
	"ZAR",
}

// URI ...
func URI(uri string) error {

	regex, err := regexp.Compile(dbregex)
	if err != nil {
		return err
	}
	if regex.MatchString(uri) {
		return nil
	}
	return fmt.Errorf("INVALID URI not supported")
}

// Currency ...
func Currency(currency string) error {

	for _, c := range currencies {
		if c == currency {
			return nil
		}
	}
	return fmt.Errorf("INVALID currency not supported")
}
