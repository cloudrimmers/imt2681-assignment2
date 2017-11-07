package tool

import (
	"log"
	"testing"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.MustLoad("../../.env")
	log.Println("!!! GOTENV !!!")
}

func TestValidateWebhook(t *testing.T) {

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
		// Fails
		{
			WebhookURL:      "http//wwww.google.com",
			BaseCurrency:    "EUR",
			TargetCurrency:  "KOK",
			CurrentRate:     1.5,
			MinTriggerValue: -1.0,
			MaxTriggerValue: 10000.0,
		},
		{
			WebhookURL:      "http://wwww.google.com",
			BaseCurrency:    "EUR",
			TargetCurrency:  "KOKK",
			CurrentRate:     1.5,
			MinTriggerValue: 1.0,
			MaxTriggerValue: 10000.0,
		},
		{
			WebhookURL:      "http://wwww.google.com",
			BaseCurrency:    "HEI",
			TargetCurrency:  "NOK",
			CurrentRate:     1.5,
			MinTriggerValue: 1.0,
			MaxTriggerValue: 10000.0,
		},
		{
			WebhookURL:      "http://wwww.google.com",
			BaseCurrency:    "EUR",
			TargetCurrency:  "NOK",
			CurrentRate:     1.5,
			MinTriggerValue: 50.0,
			MaxTriggerValue: 20.0,
		},
	}

	// Correct
	err := ValidateWebhook(&hooks[0])
	if err != nil {
		t.Error("ValidateWebhook error 1: ", err.Error())
	}

	// Fails
	err = ValidateWebhook(&hooks[1])
	if err == nil {
		t.Error("ValidateWebhook error 2: ", err.Error())
	}

	err = ValidateWebhook(&hooks[2])
	if err == nil {
		t.Error("ValidateWebhook error 3: ", err.Error())
	}

	err = ValidateWebhook(&hooks[3])
	if err == nil {
		t.Error("ValidateWebhook error 4: ", err.Error())
	}

	err = ValidateWebhook(&hooks[4])
	if err == nil {
		t.Error("ValidateWebhook error 5: ", err.Error())
	}
}
