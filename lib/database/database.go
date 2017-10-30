package database

import (
	"os"

	"gopkg.in/mgo.v2"
)

var mongoURI string = os.Getenv("MONGODB_URI")
var mongoDB string = os.Getenv("MONGODB_NAME")

func Get() (*mgo.Database, error) {
	session, err := mgo.Dial(mongoURI)
	database := session.DB(mongoDB)
	return database, err
}
