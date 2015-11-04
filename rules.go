package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ARCHIVE_DIR_NAME = "old"

	SKIPPED = ""
)

var (
	errNextFilter      = errors.New("next filter")
	errStopPropagation = errors.New("stop")
	errProcessed       = errors.New("processed")
)

type ruleFilter func(string, string, os.FileInfo) (string, error)

func filterWhitelist(root, path string, info os.FileInfo) (string, error) {
	if strings.HasPrefix(info.Name(), ARCHIVE_DIR_NAME) {
		return SKIPPED, nil
	}

	if strings.HasPrefix(info.Name(), ".") {
		return SKIPPED, nil
	}

	return SKIPPED, errNextFilter
}

func filterOlderThanOneDay(root, path string, info os.FileInfo) (string, error) {
	now := time.Now()

	if now.Sub(info.ModTime()) > time.Duration(24)*time.Hour {
		return filepath.Join(root, ARCHIVE_DIR_NAME), errProcessed
	}

	return SKIPPED, errNextFilter
}

func filterNoop(string, string, os.FileInfo) (string, error) {
	return SKIPPED, errStopPropagation
}
