package reflectUtil

import (
	"strconv"
	"testing"
)

func TestGetCallerName(t *testing.T) {

	tests := map[int]string{
		-1: "", //nonexistant
		0:  "runtime.Callers",
		1:  "reflectUtil.getCallerName",
		2:  "reflectUtil.TestGetCallerName",
		3:  "testing.tRunner",
		4:  "runtime.goexit",
		5:  "", //nonexistant
	}
	//NOTE: golangs testing framework adds another layer
	//
	for i, want := range tests {
		funcName := getCallerName(i)
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if funcName != want {
				t.Errorf("getCallerName(%d)=%v , want: %v", i, funcName, want)
			}
		})

	}
}
