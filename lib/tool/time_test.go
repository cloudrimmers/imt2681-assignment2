package tool

import (
	"testing"
	"time"
)

func TestUntilTomorrow(t *testing.T) {

	d := UntilTomorrow()
	if d >= (time.Hour * 24) {
		t.Error("time.duration is >= to 24 hours")
	}

	str := Todaystamp()
	if len(str) != 10 {
		t.Error("Malformed datestamp, to long or to short")
	}

	str1 := Daystamp(0)
	str2 := Daystamp(1)
	if len(str1) != 10 && len(str2) != 10 && str1 != str2 {
		t.Error("Malformed datestamp, to long or to short")
	}
}
