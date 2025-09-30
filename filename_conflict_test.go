package main

import (
	"os"
	"testing"
)

func TestGetOutputFileName(t *testing.T) {
	// Create a temporary file to simulate existing output file
	tmpFile, err := os.CreateTemp("", "test.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Rename to .json to simulate output file
	outputFile := tmpFile.Name() + ".json"
	err = os.Rename(tmpFile.Name(), outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile)

	// The input file name without extension
	inputFile := outputFile[:len(outputFile)-6] + ".md"

	// First call should return outputFile since it exists
	got := getOutputFileName(inputFile)
	if got == outputFile {
		t.Errorf("Expected different filename due to conflict, got %s", got)
	}

	// Create the first conflict file
	conflictFile1 := outputFile[:len(outputFile)-6] + "-1.json"
	f1, err := os.Create(conflictFile1)
	if err != nil {
		t.Fatal(err)
	}
	f1.Close()
	defer os.Remove(conflictFile1)

	got2 := getOutputFileName(inputFile)
	if got2 == outputFile || got2 == conflictFile1 {
		t.Errorf("Expected new filename different from %s and %s, got %s", outputFile, conflictFile1, got2)
	}
}
