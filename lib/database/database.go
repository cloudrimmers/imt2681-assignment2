package database

import (
	"log"

	"gopkg.in/mgo.v2"
)

var err error

// Mongo ...
type Mongo struct {
	Name    string
	URI     string
	session *mgo.Session
}

// Open ...
func (mongo *Mongo) Open() (*mgo.Database, error) {

	mongo.session, err = mgo.Dial(mongo.URI)
	if err != nil {
		return nil, err
	}
	mongo.session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it
	return mongo.session.DB(mongo.Name), nil
}

// Close ...
func (mongo *Mongo) Close() {
	mongo.session.Close()
}

// OpenC - opens a collection
func (mongo *Mongo) OpenC(cName string) (*mgo.Collection, error) {

	mongo.session, err = mgo.Dial(mongo.URI)
	if err != nil {
		return nil, err
	}
	mongo.session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it

	return mongo.session.DB(mongo.Name).C(cName), nil
}

// EnsureIndex ...
func (mongo *Mongo) EnsureIndex(cName string, keys []string) {
	// 1. Open collection
	log.Println("Ensuring unique fixer index...")
	collection, err := mongo.OpenC(cName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer mongo.Close()
	// 2. Add index
	err = collection.EnsureIndex(mgo.Index{
		Key:      keys,
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		log.Println(err.Error())
	}
}
