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

	_, fn, _, _ := runtime.Caller(1)
	envpath := filepath.Join(filepath.Dir(fn), ".env")
	log.Println("env: ", arg[1], "\nenvpath: ", envpath)

	switch arg[1] {
	case envDocker:
		break
	case envHeroku:
		break
	case envLocal:
		if err := gotenv.Load(envpath); err != nil {
			return errorInvalidPath
		}
	default:
		return errorInvalidArgument
	}

	return nil
}
