package tool

import (
	"errors"
	"net/url"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
)

func validateURI(URI string) error {
	_, err := url.ParseRequestURI(URI)
	return err
}

func validateCurrency(currency string, currencies []string) error {

	for _, c := range currencies {
		if c == currency {
			return nil
		}
	}
	return errors.New("currency not supported")
}

func validateTriggerValue(value float64, min float64, max float64) error {
	if value >= min && value <= max {
		return nil
	}
	return errors.New("trigger out of bounds")
}

// ValidateWebhook does just that
func ValidateWebhook(hook *types.Webhook, conf *types.WebConfig) error {

	var err error
	if err = validateURI(hook.WebhookURL); err != nil {
		return err
	}

	if err = validateCurrency(hook.BaseCurrency, conf.Currencies); err != nil {
		return err
	}

	if err = validateCurrency(hook.TargetCurrency, conf.Currencies); err != nil {
		return err
	}

	if err = validateTriggerValue(hook.MinTriggerValue, conf.MinTrigger, conf.MaxTrigger); err != nil {
		return err
	}

	return validateTriggerValue(hook.MaxTriggerValue, conf.MinTrigger, conf.MaxTrigger)
}
