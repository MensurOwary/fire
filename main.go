package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type arguments struct {
	from, to string
}

func args() arguments {
	args := os.Args

	if len(args) != 3 {
		exit("Only 2 arguments expected, provided %d", len(args)-1)
	}

	from := args[1]
	to := args[2]
	return arguments{
		from, to,
	}
}

func main() {
	args := args()

	fromRegex := regexp.MustCompile(args.from)

	wd, err := os.Getwd()
	panicIfNeeded(err)

	dir, err := os.Open(wd)
	panicIfNeeded(err)

	contents, err := ioutil.ReadDir(dir.Name())
	panicIfNeeded(err)

	i := 1

	for _, file := range contents {
		name := file.Name()
		if file.IsDir() || !fromRegex.MatchString(name) {
			continue
		}
		renameTo := fromRegex.ReplaceAllString(name, args.to)
		if strings.Contains(renameTo, "#{o}") {
			renameTo = strings.Replace(renameTo, "#{o}", strconv.Itoa(i), 1)
			i = i + 1
		}

		if err = os.Rename(name, renameTo); err != nil {
			panicIfNeeded(err)
		}
	}
}

func panicIfNeeded(err error) {
	exit("Error occurred : %s", err.Error())
}

func exit(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}
