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
	WebhookURL      string        `bson:"webhookURL"`
	BaseCurrency    string        `bson:"baseCurrency"`
	TargetCurrency  string        `bson:"targetCurrency"`
	CurrentRate     float64       `bson:"currentRate,omitempty"`
	MinTriggerValue float64       `bson:"minTriggerValue"`
	MaxTriggerValue float64       `bson:"maxTriggerValue"`
}

// WithinBounds ...
func (hook *Webhook) WithinBounds() bool {
	return hook.CurrentRate <= hook.MaxTriggerValue && hook.MinTriggerValue >= hook.CurrentRate
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
