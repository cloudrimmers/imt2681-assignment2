package reflectUtil

import (
	"path/filepath"
	"runtime"
)

func GetCallerName() string {
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
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
