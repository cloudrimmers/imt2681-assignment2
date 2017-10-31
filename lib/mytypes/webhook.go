package mytypes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
	"gopkg.in/mgo.v2/bson"
)

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
	ID              bson.ObjectId `json:"id" bson:"_id,omitempty"`
	WebhookURL      string        `json:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency"`
	MinTriggerValue float64       `json:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue"`
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
	ID              bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BaseCurrency    string        `json:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency"`
	CurrentRate     float64       `json:"currentRate"`
	MinTriggerValue float64       `json:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue"`
}

// Trigger ...
func (hook *mytypes.WebhookIn) Trigger() {

	data, _ := json.Marshal(hook) // @TODO this should actually be a webhookOut structure
	req, _ := http.NewRequest("POST", hook.WebhookURL, &data)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error triggering webhook..", err.Error())
	}
	defer resp.Body.Close()
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
