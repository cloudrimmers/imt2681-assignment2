package tool

import (
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FireWebhooks ...
func FireWebhooks(collectionWebhook *mgo.Collection, collectionFixer *mgo.Collection) {
	hooks := []types.Webhook{}
	collectionWebhook.Find(nil).All(&hooks)

	for _, hook := range hooks {

		fixer := types.FixerIn{}
		err := collectionFixer.Find(bson.M{"base": hook.BaseCurrency}).Sort("-date").One(&fixer)
		if err != nil {
			panic("No hook.baseCurrency found in the database!")
		}
		hook.CurrentRate = fixer.Rates[hook.TargetCurrency]
		go hook.Trigger()

	}
}
