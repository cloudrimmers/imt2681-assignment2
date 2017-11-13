package dialogFlow

import (
	"strconv"
	"testing"
)

func TestNewQuery(t *testing.T) {
	//SETUP
	tests := []string{
		"Convert 200 bucks to nok",
		"1¥ as $",
	}
	for i, test := range tests {

		qry := newQuery(test)
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test != qry.Query {
				wanted := query{
					Query:     test,
					Contexts:  qry.Contexts,
					SessionID: qry.SessionID,
				}
				t.Errorf("newQuery() = %v want %v", qry, wanted)
			}
			if len(qry.SessionID) <= 0 {
				wanted := query{
					Query:     test,
					Contexts:  qry.Contexts,
					SessionID: "NOT EMPTY",
				}
				t.Errorf("newQuery() = %v want %v", qry, wanted)
			}
		})
	}
	//Teardown
}
