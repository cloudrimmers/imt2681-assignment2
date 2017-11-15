package environment

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/subosito/gotenv"
)

var (
	envDocker            = "docker"
	envHeroku            = "heroku"
	envLocal             = "local"
	errorInvalidArgument = fmt.Errorf("ERROR Invalid command-line argument")
	errorInvalidPath     = fmt.Errorf("ERROR Invalid envpath")
)

// Load ...
func Load(arg []string) error {
	debugprint := func(envname string, envpath string) {
		log.Println("env name: ", envname)
		log.Println("envpath:  ", envpath)
	}
	_, fn, _, _ := runtime.Caller(1)
	envpath := filepath.Join(filepath.Dir(fn), ".env")

	switch arg[1] {
	case envDocker:
		debugprint(arg[1], "using existing environment")
	case envHeroku:
		debugprint(arg[1], "using existing environment")
	case envLocal:
		debugprint(arg[1], envpath)
		if err := gotenv.Load(envpath); err != nil {
			return errorInvalidPath
		}
	default:
		debugprint(arg[1], "wrong env name")
		return errorInvalidArgument
	}

	return nil
}
