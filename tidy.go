package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	ruleFilters = []ruleFilter{
		filterWhitelist,
		filterOlderThanOneDay,
		filterNoop,
	}
)

func main() {
	path := mustGetCleanablePath()
	fileGroups := mustFilterFiles(path)
	mustTidyFiles(fileGroups)
}

func mustGetCleanablePath() (path string) {
	var err error

	flag.Parse()

	if len(flag.Args()) < 1 {
		path = "."
	} else {
		path = flag.Args()[0]
	}

	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory\n", path)
		os.Exit(1)
	}

	return
}

func mustFilterFiles(root string) map[string]string {
	groups := make(map[string]string)

	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while processing %s %q\n", path, err)
			return nil
		}

		for _, ruleFilter := range ruleFilters {
			to, err := ruleFilter(root, path, info)

			if err == errStopPropagation {
				return nil
			}
			if err == errNextFilter {
				continue
			}
			if err == errProcessed {
				groups[path] = to
				return nil
			}
		}

		return nil
	})
	if walkErr != nil {
		panic(walkErr)
	}

	return groups
}

func mustTidyFiles(fileGroups map[string]string) {
	var (
		err       error
		processed = 0
	)

	// Create directories.
	for _, toDir := range fileGroups {
		log.Printf("toDir: %s", toDir)
		if err := os.MkdirAll(toDir, os.ModeDir|0755); err != nil {
			panic(err)
		}
	}

	// Move files.
	for fileName, toDir := range fileGroups {
		newFileName := filepath.Join(toDir, filepath.Base(fileName))
		if newFileName == fileName {
			continue
		}
		if err = os.Rename(fileName, newFileName); err != nil {
			panic(err)
		} else {
			processed = processed + 1
			fmt.Fprintf(os.Stderr, "mv %s -> %s\n", fileName, newFileName)
		}
	}
	fmt.Fprintf(os.Stderr, "processed %d files\n", processed)
}
