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

var mongoURI string = os.Getenv("MONGODB_URI")
var mongoDB string = os.Getenv("MONGODB_NAME")
var session *mgo.Session
var err error

// EnsureFixerIndex ...
func EnsureFixerIndex(collectionFixer string) {

	// 1. Open database
	log.Println("Ensuring unique fixer index...")
	db, err := Open()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer Close()

	index := mgo.Index{
		Key:      []string{"date"},
		Unique:   true,
		DropDups: true,
	}

	err = db.C(collectionFixer).EnsureIndex(index)
	if err != nil {
		log.Println(err.Error())
	}
}

// SeedFixer ...
func SeedFixer(collectionFixer string) {
	// 1. Open database
	db, err := Open()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer Close()

	// 2. Read from file
	basepath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	data, err := ioutil.ReadFile(basepath + "/seedfixer.json")
	log.Println("loading seed data from ", basepath+"/seedfixer.json")
	fixerData := []types.FixerIn{}

	if err != nil {
		panic(err.Error())
	}
	if err = json.Unmarshal(data, &fixerData); err != nil {
		panic(err.Error())
	}
	// 3. Insert to database
	// db.C(collectionFixer).DropCollection()
	for _, fixer := range fixerData {
		if err = db.C(collectionFixer).Insert(fixer); err != nil {
			log.Println("Unable to db.Insert seed")
		}
	}
}

// Open ...
func Open() (*mgo.Database, error) {
	database := &mgo.Database{}
	session, err = mgo.Dial(mongoURI)

	session.SetMode(mgo.Monotonic, true) // @note not sure what this is, but many people use it

	if err == nil {
		database = session.DB(mongoDB)
	}
	return database, err
}

// Close ...
func Close() {
	session.Close()
}
