package main

import (
	"os"
	"strings"
	"testing"
)

func TestMissingInputFile(t *testing.T) {
	exitCode := parseArgs([]string{}, &strings.Builder{})
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 for missing args, got %d", exitCode)
	}
}

func TestNonExistentInputFile(t *testing.T) {
	if checkFileExists("nonexistentfile.md") {
		t.Error("Expected nonexistent file to return false")
	}
}

func TestOutputFileNameConflict(t *testing.T) {
	// Create a temporary file to simulate existing output file
	tmpFile, err := os.CreateTemp("", "test.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	outputFile := tmpFile.Name() + ".json"
	err = os.Rename(tmpFile.Name(), outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile)

	got := getOutputFileName(tmpFile.Name())
	if got == outputFile {
		t.Errorf("Expected different filename due to conflict, got %s", got)
	}
}
