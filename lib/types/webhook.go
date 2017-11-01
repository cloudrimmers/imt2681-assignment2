package types

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	errorCurrency     = errors.New("invalid currency format")
	errorTriggerValue = errors.New("invalid trigger value")
)

// Webhook central datastructure for the service
/* Example:
{
	"webhookURL": "http://remoteUrl:8080/randomWebhookPath",
	"baseCurrency": "EUR",
	"targetCurrency": "NOK",
	"minTriggerValue": 1.50,
	"maxTriggerValue": 2.55
}
*/
type Webhook struct {
	ID              bson.ObjectId `json:"id" bson:"_id,omitempty"`
	WebhookURL      string        `json:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency"`
	CurrentRate     float64       `json:"currentRate,omitempty" bson:",omitempty"`
	MinTriggerValue float64       `json:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue"`
}

// CurrencyIn ...
/* Example:
{
	"baseCurrency": "EUR",
	"targetCurrency": "NOK",
}
*/
type CurrencyIn struct {
	BaseCurrency   string `json:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency"`
}
