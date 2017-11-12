package timetool

import (
	"testing"
	"time"
)

func TestUntilTomorrow(t *testing.T) {

	d := UntilTomorrow()
	if d >= (time.Hour * 24) {
		t.Error("time.duration is >= to 24 hours")
	}

}
