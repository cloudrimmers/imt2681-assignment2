package tool

import (
	"github.com/Arxcis/imt2681-assignment2/lib/types"
	mgo "gopkg.in/mgo.v2"
)

// FireWebhooks ...
func FireWebhooks(collectionWebhook *mgo.Collection, collectionFixer *mgo.Collection) {
	hooks := []types.Webhook{}
	collectionWebhook.Find(nil).All(&hooks)

	for _, hook := range hooks {
		go hook.Trigger(collectionFixer)
	}
}
