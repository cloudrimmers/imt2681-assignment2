package database

import (
	"fmt"
	"testing"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.MustLoad("../../.env")

}
func TestOpen(t *testing.T) {

	mongo := Mongo{Name: "test", URI: "127.0.0.1:33017", Session: nil}

	// Test 1 - Open() - SUCCESS CASE
	_, err := mongo.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()

	// Test 2 - FAIL CASE

	mongo = Mongo{Name: "test", URI: "trolololol", Session: nil}

	_, err = mongo.Open()
	if err == nil {
		defer mongo.Close()
		t.Error(fmt.Errorf("ERROR Should not have been able to access database with URI: ", mongo.URI))
	}

	fmt.Println("opne(3)")

	mongo = Mongo{Name: "test", URI: "127.0.0.1:33017", Session: nil}
	// Test 2
	_, err = mongo.OpenC("test")
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()
}

func TestEnsureIndex(t *testing.T) {

	mongo := Mongo{Name: "test", URI: "mongodb://127.0.0.1:33017", Session: nil}
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
