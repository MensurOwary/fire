package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetCurrentDirContents() []os.FileInfo {
	wd, err := os.Getwd()
	Exit(err)

	dir, err := os.Open(wd)
	Exit(err)

	contents, err := ioutil.ReadDir(dir.Name())
	Exit(err)

	return contents
}

func Exit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occurred : %s\n", err.Error())
		os.Exit(2)
	}
}
