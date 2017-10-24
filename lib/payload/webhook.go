package payload

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
	BaseCurrency    string
	TargetCurrency  string
	CurrentRate     float64
	MinTriggerValue float64
	MaxTriggerValue float64
}
