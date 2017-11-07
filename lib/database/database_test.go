package database

import (
	"testing"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.MustLoad("../../.env")

}
func TestDatabase(t *testing.T) {

	_, err := Open()
	if err != nil {
		t.Error(err.Error())
	}

	Close()
}
