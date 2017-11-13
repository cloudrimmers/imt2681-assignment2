package database

import (
	"fmt"
	"log"

	"github.com/Arxcis/imt2681-assignment2/lib/validate"
	"gopkg.in/mgo.v2"
)

var err error

// Mongo ...
type Mongo struct {
	Name    string
	URI     string
	Session *mgo.Session
}

// Open ...
func (mongo *Mongo) Open() (*mgo.Database, error) {

	validate.URI()

	mongo.Session, err = mgo.Dial(mongo.URI)

	if err != nil {
		fmt.Println("awerwerwererewrewrw", err.Error())
		return nil, err
	}
	log.Println("Database connection established...")
	mongo.Session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it
	return mongo.Session.DB(mongo.Name), nil
}

// Close ...
func (mongo *Mongo) Close() {
	mongo.Session.Close()
}

// OpenC - opens a collection
func (mongo *Mongo) OpenC(cName string) (*mgo.Collection, error) {

	mongo.Session, err = mgo.Dial(mongo.URI)
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established")
	mongo.Session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it

	return mongo.Session.DB(mongo.Name).C(cName), nil
}

// EnsureIndex ...
func (mongo *Mongo) EnsureIndex(cName string, keys []string) error {
	// 1. Open collection
	log.Println("Ensuring unique " + cName + " index")
	collection, err := mongo.OpenC(cName)
	if err != nil {
		return err
	}
	defer mongo.Close()
	log.Println("Success ensuring " + cName + " index")
	// 2. Add index
	err = collection.EnsureIndex(mgo.Index{
		Key:      keys,
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		return err
	}
	return nil
}
