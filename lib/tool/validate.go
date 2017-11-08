package tool

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
)

var currencies []string

// Load the settings file to configure settings
func loadCurrencies(filepath string) []string {

	//basepath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	data, err := ioutil.ReadFile(filepath)
	log.Println("Loading config from : ", filepath)

	if err != nil {
		log.Println("Panic upcomming.....")
		panic(err.Error())

	}
	var currencies []string
	if err = json.Unmarshal(data, &currencies); err != nil {
		panic(err.Error())
	}
	//	log.Println("Validation settings file: ", v)
	return currencies
}

func validateURI(URI string) error {
	_, err := url.ParseRequestURI(URI)
	return err
}

func validateCurrency(currency string) error {

	for _, c := range currencies {

		if c == currency {
			return nil
		}
	}
	return errors.New("currency not supported")
}

func validateTriggerValue(min float64, max float64) error {
	if min < max && min >= 0.0 && max > 0.0 {
		return nil
	}
	return errors.New("trigger out of bounds")
}

// ValidateWebhook does just that
func ValidateWebhook(hook *types.Webhook) error {

	var err error
	if err = validateURI(hook.WebhookURL); err != nil {
		return err
	}

	if err = validateCurrency(hook.BaseCurrency); err != nil {
		return err
	}

	if err = validateCurrency(hook.TargetCurrency); err != nil {
		return err
	}

	return validateTriggerValue(hook.MinTriggerValue, hook.MaxTriggerValue)
}
