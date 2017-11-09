package validate

import (
	"testing"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
)

func TestWebhook(t *testing.T) {

	hooks := []types.Webhook{
		// Correct
		{
			WebhookURL:      "http://wwww.google.com",
			BaseCurrency:    "EUR",
			TargetCurrency:  "NOK",
			CurrentRate:     1.5,
			MinTriggerValue: 1.4,
			MaxTriggerValue: 1.6,
		},
	}
	currency := []string{"EUR", "NOK", "USD", "SEK"}

	// Correct
	err := Webhook(&hooks[0], currency)
	if err != nil {
		t.Error("Webhook error 1: ", err.Error())
	}
}
