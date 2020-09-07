package arg

import (
	"errors"
	"github.com/mensurowary/fire/util"
	"os"
	"strconv"
	s "strings"
)

// Arguments represents the application arguments
type Arguments struct {
	From, To                string
	IncludeDir, IncludeFile bool
	Async                   bool
}

// Args function constructs the input arguments from cli arguments
func Args() Arguments {
	osArgs := os.Args

	l := len(osArgs)

	if !(l == 3 || l == 4) {
		util.Exit(errors.New("only 2 arguments expected, provided " + strconv.Itoa(len(osArgs)-1)))
	}

	if l == 3 {
		return defaultMode(osArgs)
	}

	var includeDir, includeFile, async bool

	params := osArgs[1]
	if s.Index(params, "-") != 0 {
		util.Exit(errors.New("parameters expected to have '-' at the beginning"))
	}

	// evaluating commands
	for _, p := range params[1:] {
		switch string(p) {
		case "f":
			includeFile = true
			break
		case "d":
			includeDir = true
		case "a":
			async = true
		}
	}
	return Arguments{
		osArgs[2],
		osArgs[3],
		includeDir,
		includeFile,
		async,
	}
}

func defaultMode(osArgs []string) Arguments {
	return Arguments{
		osArgs[1], osArgs[2],
		false, true, false,
	}
}
