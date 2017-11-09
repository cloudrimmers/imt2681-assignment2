package database

import (
	"testing"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.MustLoad("../../.env")

}
func TestDatabase(t *testing.T) {

	mongo := Mongo{Name: "test", URI: "127.0.0.1:33017", Session: nil}

	// Test 1
	_, err := mongo.Open()
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()

	// Test 2
	_, err = mongo.OpenC("test")
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()

	// Test 3
	err = mongo.EnsureIndex("test", []string{"id"})
	if err != nil {
		t.Error(err.Error())
	}
	mongo.Close()
}
