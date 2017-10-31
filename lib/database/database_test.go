package database

import "testing"

func TestDatabase(t *testing.T) {

	_, err := Open()
	if err != nil {
		t.Error(err.Error())
	}

	Close()
}
