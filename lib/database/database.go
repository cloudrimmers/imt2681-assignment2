package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Arxcis/imt2681-assignment2/lib/types"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session
var err error

// Open ...
func Open() (*mgo.Database, error) {

	var mongoURI = os.Getenv("MONGODB_URI")
	var mongoDB = os.Getenv("MONGODB_NAME")

	session, err = mgo.Dial(mongoURI)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it
	return session.DB(mongoDB), nil
}

// OpenWebhook ...
func OpenWebhook() (*mgo.Collection, error) {

	var mongoURI = os.Getenv("MONGODB_URI")
	var mongoDB = os.Getenv("MONGODB_NAME")
	var collectionWebhook = os.Getenv("COLLECTION_WEBHOOK")

	session, err = mgo.Dial(mongoURI)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it

	return session.DB(mongoDB).C(collectionWebhook), nil
}

// OpenFixer ...
func OpenFixer() (*mgo.Collection, error) {

	var mongoURI = os.Getenv("MONGODB_URI")
	var mongoDB = os.Getenv("MONGODB_NAME")
	var collectionFixer = os.Getenv("COLLECTION_FIXER")

	session, err = mgo.Dial(mongoURI)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it

	return session.DB(mongoDB).C(collectionFixer), nil
}

// Close ...
func Close() {
	session.Close()
}

// EnsureFixerIndex ...
func EnsureFixerIndex(collectionFixer string) {

	// 1. Open collection
	log.Println("Ensuring unique fixer index...")
	dbfixer, err := OpenFixer()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer Close()

	// 2. Add index
	err = dbfixer.EnsureIndex(mgo.Index{
		Key:      []string{"date"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

// SeedFixer ...
func SeedFixer() {

	var seedPath = os.Getenv("SEED_PATH")

	// 1. Read from file
	basepath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fullpath := basepath + seedPath
	data, err := ioutil.ReadFile(fullpath)
	log.Println("loading seed data from ", fullpath)
	fixerData := []types.FixerIn{}

	if err != nil {
		panic(err.Error())
	}

	// 2. Unmarshal
	if err = json.Unmarshal(data, &fixerData); err != nil {
		panic(err.Error())
	}

	// 3. Open collection
	cfixer, err := OpenFixer()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer Close()

	// 4. Insert to database
	// cfixer.DropCollection()
	for _, fixer := range fixerData {
		if err = cfixer.Insert(fixer); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}
