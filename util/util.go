package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

// GetCurrentDirContents function returns an array of files and directories
// that reside in the current directory
// The function exits with the status code 2 failure and the appropriate message
func GetCurrentDirContents() []os.FileInfo {
	wd, err := os.Getwd()
	Exit(err)

	dir, err := os.Open(wd)
	Exit(err)

	contents, err := ioutil.ReadDir(dir.Name())
	Exit(err)

	return contents
}

// Exit method is the utility method to exit with the given error and status code 2
func Exit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occurred : %s\n", err.Error())
		os.Exit(2)
	}
}
