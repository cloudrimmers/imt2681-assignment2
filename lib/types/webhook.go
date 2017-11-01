package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"gopkg.in/mgo.v2"

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

// Trigger triggers a request to the webhook to the WebhookURL
func (hook *Webhook) Trigger(collectionFixer *mgo.Collection) {

	fixer := types.FixerIn{}
	err := collectionFixer.Find(bson.M{"base": hook.BaseCurrency}).Sort("-date").One(&fixer)
	if err != nil {
		panic("No hook.baseCurrency found in the database!")
	}
	hook.CurrentRate = fixer.Rates[hook.TargetCurrency]

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
