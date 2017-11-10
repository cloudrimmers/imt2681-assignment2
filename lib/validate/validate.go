package validate

import (
	"errors"
	"net/url"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
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
	return errors.New("currency not supported")
}

// TriggerValue ...
func TriggerValue(min float64, max float64) error {
	if min < max && min >= 0.0 && max > 0.0 {
		return nil
	}
	return errors.New("trigger out of bounds")
}

// NewWebhook ...
func NewWebhook(hook *types.Webhook, currency []string) error {

	var err error
	if err = URI(hook.WebhookURL); err != nil {
		return err
	}

	if err = Currency(hook.BaseCurrency, currency); err != nil {
		return err
	}

	if err = Currency(hook.TargetCurrency, currency); err != nil {
		return err
	}

	return TriggerValue(hook.MinTriggerValue, hook.MaxTriggerValue)
}
