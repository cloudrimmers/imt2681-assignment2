package reflectUtil

import (
	"path/filepath"
	"runtime"
)

// REFERENCE: https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang

//GetCallerName - returns the name of the function the used
func GetCallerName() string {
	return getCallerName(3)
}

func GetCallerNameInTest() string {
	return getCallerName(2)
}

func getCallerName(depth int) string {
	fpcs := make([]uintptr, 1)

	// skip 2 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(depth, fpcs)
	if n == 0 {
		return "" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return ""
	}
	_, fileName := filepath.Split(fun.Name())
	return fileName
}
