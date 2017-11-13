package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

var mongoURI string

func init() {
	gotenv.MustLoad("../../.env")
	mongoURI = os.Getenv("MONGODB_URI")
}
func TestOpen(t *testing.T) {

	mongo := Mongo{Name: "test", URI: mongoURI, Session: nil}

	// Test 1 - Open() - SUCCESS CASE
	_, err = mongo.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()

	// Test 2 - FAIL CASE

	mongo = Mongo{Name: "test", URI: "trolololol", Session: nil}

	_, err = mongo.Open()
	if err == nil {
		defer mongo.Close()
		t.Error("ERROR Should not have been able to access database with URI: ", mongo.URI)
	}

	fmt.Println("opne(3)")

	mongo = Mongo{Name: "test", URI: mongoURI, Session: nil}
	// Test 2
	_, err = mongo.OpenC("test")
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()
}

func TestEnsureIndex(t *testing.T) {

	mongo := Mongo{Name: "test", URI: mongoURI, Session: nil}
	err = mongo.EnsureIndex("test", []string{"id"})
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()

	// FAIL TEST
	mongo = Mongo{Name: "test", URI: "127.", Session: nil}
	err = mongo.EnsureIndex("test", []string{"id"})
	if err == nil {
		defer mongo.Close()
		t.Error(err.Error())
	}
}
