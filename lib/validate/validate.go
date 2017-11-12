package validate

import (
	"fmt"
	"net/url"
)

// URI ...
func URI(uri string) error {
	_, err := url.ParseRequestURI(uri)
	return err
}

// Currency ...
func Currency(currency string, currencyArray []string) error {

	for _, c := range currencyArray {
		if c == currency {
			return nil
		}
	}
	return fmt.Errorf("currency not supported")
}
