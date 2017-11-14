package validate

import (
	"testing"
)

func TestURI(t *testing.T) {

	table := map[string]bool{
		"mongodb://localhost": true,
		"127.0.0.1":           true,
		"mongodb://heroku_xxxxx:xxxxxxx@ds239965.mlab.com:33333/xxxxxx": true,
		"127.0":    false,
		"dfklsdfk": false,
	}

	for uri, want := range table {

		t.Run(uri, func(t *testing.T) {

			err := URI(uri)
			if err == nil {
				if !want {
					t.Error("URI  matched but we didn't want it to.")
				}
			} else {
				if want {
					t.Error("URI  did not match but we wanted it to.")
				}
			}
		})
	}

	// fiddle with the regexp
	dbregex = "^^^^^)¤¤¤$$$$$$$$"
	if err := URI("localhost"); err == nil {
		t.Error(err.Error())
	}

}

func TestCurrencies(t *testing.T) {

	table := map[string]bool{ //Bool determines rather the test is suppose to fail.
		"NOK":    true,
		"EUR":    true,
		"blalba": false,
		"_dsfdf": false,
		"kok":    false,
		"":       false,
	}

	for cur, want := range table {

		t.Run(cur, func(t *testing.T) {

			err := Currency(cur)
			if err == nil {
				if !want {
					t.Error("CURRENCY matched but we didn't want it to.")
				}
			} else {
				if want {
					t.Error("CURRENCY did not match but we wanted it to.")
				}
			}
		})

	}
}
