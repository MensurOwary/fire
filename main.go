package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	r "regexp"
	"strconv"
	s "strings"
)

var orderSymbol = "#{o}"

type arguments struct {
	from, to string
}

func args() arguments {
	args := os.Args

	if len(args) != 3 {
		exit(errors.New("Only 2 arguments expected, provided " + strconv.Itoa(len(args)-1)))
	}

	from := args[1]
	to := args[2]
	return arguments{
		from, to,
	}
}

func main() {
	args := args()

	fromRegex := r.MustCompile(args.from)

	i := 1

	for _, file := range getCurrentDirContents() {
		from := file.Name()
		if file.IsDir() || !fromRegex.MatchString(from) {
			continue
		}
		to := fromRegex.ReplaceAllString(from, args.to)
		if s.Contains(to, orderSymbol) {
			to = s.Replace(to, orderSymbol, strconv.Itoa(i), 1)
			i = i + 1
		}

		if err := os.Rename(from, to); err != nil {
			exit(err)
		}
	}
}

func getCurrentDirContents() []os.FileInfo {
	wd, err := os.Getwd()
	exit(err)

	dir, err := os.Open(wd)
	exit(err)

	contents, err := ioutil.ReadDir(dir.Name())
	exit(err)

	return contents
}

func exit(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occurred : %s\n", err.Error())
	os.Exit(2)
}
