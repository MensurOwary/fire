package main

import (
	"github.com/mensurowary/fire/arg"
	"github.com/mensurowary/fire/util"
	"github.com/mensurowary/fire/wildcard"
	"log"
	"os"
	r "regexp"
	"sync"
)

type pair struct {
	from, to string
}

func run(args arg.Arguments) {
	pairs := makeNewNames(args)
	rename(pairs, args)
}

func makeNewNames(args arg.Arguments) []pair {
	ord := wildcard.MakeOrdering(args)
	fromRegex := r.MustCompile(args.From)

	var conversion []pair

	for _, file := range util.GetCurrentDirContents() {
		from := file.Name()
		if shouldSkip(args, file) || !fromRegex.MatchString(from) {
			continue
		}
		to := fromRegex.ReplaceAllString(from, args.To)
		to = ord.ReplaceIfNeeded(to)

		conversion = append(conversion, pair{from, to})
	}

	if conversion == nil {
		conversion = []pair{}
	}

	return conversion
}

func rename(pairs []pair, args arg.Arguments) {
	if args.Async {
		log.Println("processing in an async way")
		renameAsync(pairs)
	} else {
		renameSync(pairs)
	}
}

// renaming all in one loop
func renameSync(pairs []pair) {
	for _, file := range pairs {
		if err := os.Rename(file.from, file.to); err != nil {
			util.Exit(err)
		}
		log.Printf("Renamed %s => %s\n", file.from, file.to)
	}
}

// it just divides the pairs into 2 and renames them separately
func renameAsync(pairs []pair) {
	l := len(pairs)

	first := pairs[:l/2]
	second := pairs[l/2:]

	var wg sync.WaitGroup
	wg.Add(2)

	go doRename(first, &wg)
	go doRename(second, &wg)

	wg.Wait()
}

func doRename(first []pair, wg *sync.WaitGroup) {
	for _, file := range first {
		if err := os.Rename(file.from, file.to); err != nil {
			util.Exit(err)
		}
	}
	wg.Done()
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
