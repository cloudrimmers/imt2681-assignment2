package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

var (
	errorURL          = errors.New("invalid url format")
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

// Trigger triggers a request to the webhook to the WebhookURL
func (hook *Webhook) Trigger() {

	data, _ := json.Marshal(hook) // @TODO this should actually be a webhookOut structure
	reader := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", hook.WebhookURL, reader)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error triggering webhook..", err.Error())
	}
	defer resp.Body.Close()
}

// Validate the webhook data when in comes form the client
func (hook *Webhook) Validate() {

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
