package invoke

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Webhooks ...
func Webhooks(client *http.Client, cWebhook *mgo.Collection, cFixer *mgo.Collection) {
	hooks := []types.Webhook{}
	cWebhook.Find(nil).All(&hooks)

	for _, hook := range hooks {

		// 1. Find
		fixer := types.FixerIn{}
		err := cFixer.Find(bson.M{"base": hook.BaseCurrency}).Sort("-date").One(&fixer)
		if err != nil {
			panic("No hook.baseCurrency found in the database!")
		}

		// 2.5 Check if within
		hook.CurrentRate = fixer.Rates[hook.TargetCurrency]
		if !hook.WithinBounds() {

			// 2. Marshall data
			data, _ := json.Marshal(hook)
			/*discorddata := struct {
				Content string `json:"content"`
			}{
				string(dat),
			}
			data, _ := json.Marshal(&discorddata) // @TODO this should actually be a webhookOut structure
			*/
			// 3. Create a POST request
			req, err := http.NewRequest("POST", hook.WebhookURL, bytes.NewReader(data))
			if err != nil {
				log.Println("Invalid request", err.Error())
				return
			}

			req.Header.Add("content-type", "application/json")
			log.Println("Fireing webhook: ", hook)

			// 4. Do request
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error posting webhook..", err.Error())
			}
			defer resp.Body.Close()
			bytesbody, err := ioutil.ReadAll(resp.Body)
			log.Println("Response body: ", string(bytesbody))
		}
	}
	log.Println("All hooks fired and done...")
}
