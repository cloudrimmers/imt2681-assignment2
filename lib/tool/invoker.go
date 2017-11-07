package tool

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// InvokeWebhooks ...
func InvokeWebhooks(client *http.Client, collectionWebhook *mgo.Collection, collectionFixer *mgo.Collection) {
	hooks := []types.Webhook{}
	collectionWebhook.Find(nil).All(&hooks)
	for _, hook := range hooks {

		// 1. Find
		fixer := types.FixerIn{}
		err := collectionFixer.Find(bson.M{"base": hook.BaseCurrency}).Sort("-date").One(&fixer)
		if err != nil {
			panic("No hook.baseCurrency found in the database!")
		}

		// 2.5 Check if within
		hook.CurrentRate = fixer.Rates[hook.TargetCurrency]
		if !hook.WithinBounds() {

			// 2. Marshall data
			dat, _ := json.Marshal(hook)
			discorddata := struct {
				Content string `json:"content"`
			}{
				string(dat),
			}
			data, _ := json.Marshal(&discorddata) // @TODO this should actually be a webhookOut structure

			// 3. Create a POST request
			req, err := http.NewRequest("POST", hook.WebhookURL, bytes.NewReader(data))
			if err != nil {
				log.Println("Invalid request", err.Error())
				return
			}

			req.Header.Add("content-type", "application/json")
			log.Println("Fireing webhook: ", req, string(data))

			// 4. Do request
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error posting webhook..", err.Error())
			}
			defer resp.Body.Close()
			log.Println(resp)
		}
	}
	log.Println("All hooks fired and done...")
}
