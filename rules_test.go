package main

import (
	"os"
	"testing"
	"time"
)

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestFilterWhitelist(t *testing.T) {
	for _, name := range []string{".", "..", ".abc", "..abcde"} {
		dotFile := mockFileInfo{name: name}
		to, err := filterWhitelist("", "", dotFile)
		if to != SKIPPED || err != errStopPropagation {
			t.Errorf("should exclude dot files: %s %q", to, err)
		}
	}

	file := mockFileInfo{name: ARCHIVE_DIR_NAME}
	to, err := filterWhitelist("", "", file)
	if to != SKIPPED || err != errStopPropagation {
		t.Errorf("should exclude archive dir: %s %q", to, err)
	}
}

func TestFilterNoop(t *testing.T) {
	file := mockFileInfo{}
	to, err := filterNoop("", "", file)
	if to != SKIPPED || err != errStopPropagation {
		t.Errorf("should stop propagation: %s %q", to, err)
	}
}

func TestFilterOlderThanOneDay(t *testing.T) {
	archiveFile := mockFileInfo{
		modTime: time.Now().Add(time.Duration(-24) * time.Hour),
	}
	to, err := filterOlderThanOneDay("", "", archiveFile)
	if to != ARCHIVE_DIR_NAME || err != errProcessed {
		t.Errorf("should process archived file: %s %q", to, err)
	}

	freshFile := mockFileInfo{
		modTime: time.Now().Add(time.Duration(-5) * time.Hour),
	}
	to, err = filterOlderThanOneDay("", "", freshFile)
	if to != SKIPPED || err != errNextFilter {
		t.Errorf("should not process fresh file: %s %q", to, err)
	}
}
