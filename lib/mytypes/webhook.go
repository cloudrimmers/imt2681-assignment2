package mytypes

import "gopkg.in/mgo.v2/bson"

// WebhookIn ...
/* Example:
{
	"webhookURL": "http://remoteUrl:8080/randomWebhookPath",
	"baseCurrency": "EUR",
	"targetCurrency": "NOK",
	"minTriggerValue": 1.50,
	"maxTriggerValue": 2.55
}
*/
type WebhookIn struct {
	ID              bson.ObjectId
	WebhookURL      string
	BaseCurrency    string
	TargetCurrency  string
	MinTriggerValue float64
	MaxTriggerValue float64
}

// WebhookOut ...
/* Example:
{
	"baseCurrency": "EUR",
	"targetCurrency": "NOK",
	"currentRate": 2.75,
	"minTriggerValue": 1.50,
	"maxTriggerValue": 2.55
}
*/
type WebhookOut struct {
	ID              bson.ObjectId
	BaseCurrency    string
	TargetCurrency  string
	CurrentRate     float64
	MinTriggerValue float64
	MaxTriggerValue float64
}

// CurrencyIn ...
/* Example:
{
	"baseCurrency": "EUR",
	"targetCurrency": "NOK",
}
*/
type CurrencyIn struct {
	BaseCurrency   string
	TargetCurrency string
}
