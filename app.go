package main

import (
	"github.com/mensurowary/fire/arg"
	"github.com/mensurowary/fire/util"
	"github.com/mensurowary/fire/wildcard"
	"os"
	r "regexp"
)

func run(args arg.Arguments) {
	ord := wildcard.MakeOrdering(args)
	fromRegex := r.MustCompile(args.From)

	for _, file := range util.GetCurrentDirContents() {
		from := file.Name()
		if shouldSkip(args, file) || !fromRegex.MatchString(from) {
			continue
		}
		to := fromRegex.ReplaceAllString(from, args.To)
		to = ord.ReplaceIfNeeded(to)

		if err := os.Rename(from, to); err != nil {
			util.Exit(err)
		}
	}
}

func shouldSkip(args arg.Arguments, file os.FileInfo) bool {
	isDir := file.IsDir()
	return decisionTable(
		args.IncludeFile,
		args.IncludeDir,
		!isDir,
		isDir,
	)
}

// decision table to determine whether the current resource should be skipped
func decisionTable(incFile, incDir, isFile, isDir bool) bool {
	arr := [6][5]bool{
		{true, true, true, false, false},
		{true, true, false, true, false},
		{true, false, true, false, false},
		{true, false, false, true, true},
		{false, true, true, false, true},
		{false, true, false, true, false},
	}

	for _, line := range arr {
		c := line[0] == incFile &&
			line[1] == incDir &&
			line[2] == isFile &&
			line[3] == isDir
		if c {
			return line[4]
		}
	}
	return false
}
