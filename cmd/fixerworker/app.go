package main

import (
	"github.com/Arxcis/imt2681-assignment2/lib/database"
)

// App ...
type App struct {
	FixerioURI        string
	SeedFixerPath     string
	CollectionWebhook string
	CollectionFixer   string
	MongodbName string
	MongodbURI string
}

// Fixer2Mongo ...
func (APP *App) Fixer2Mongo(mongo *database.Mongo) {
	

	// 1. Connect and request to fixer.io
	resp, err := http.Get(APP.FixerioURI)
	if err != nil {
		log.Println("No connection with fixer.io: "+APP.FixerioURI+" ...", err.Error())
		return
	}

	// 2. Decode payload
	payload := &(types.FixerIn{})
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		log.Println("Could not decode resp.Body...", err.Error())
		return
	}

	// 3. Connect to DB
	dbsession, err := mongo.Open()
	if err != nil {
		log.Println("Database no connection: ", err.Error())
	}
	defer mongo.Close()

	payload.Timestamp = time.Now().String()

	// 5. Dump payload to database
	err = dbsession.C(APP.CollectionFixer).Insert(payload)
	if err != nil {
		log.Println("Error on db.Insert():\n", err.Error())

	}
	log.Println("Successfull grab of fixer.io: ", payload)

	// 6. Fire webhooks
	client := &http.Client{}
	tool.InvokeWebhooks(client, dbsession.C(APP.CollectionWebhook), dbsession.C(APP.CollectionFixer))
}
	
// SeedFixer ...
func (app *App) SeedFixer(mongo *database.Mongo) {

	// 1. Read from file
	data, err := ioutil.ReadFile(app.SeedFixerPath)
	log.Println("loading seed data from ", app.SeedFixerPath)

	if err != nil {
		panic(err.Error())
	}

	// 2. Unmarshal
	marshalled := []types.FixerIn{}
	if err = json.Unmarshal(data, &marshalled); err != nil {
		panic(err.Error())
	}

	// 3. Open collection
	collection, err := mongo.OpenC(app.CollectionFixer)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer mongo.Close()

	// 4. Insert to database
	// cfixer.DropCollection()
	for _, o := range marshalled {
		if err = collection.Insert(o); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}
