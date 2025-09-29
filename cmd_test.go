package main

import (
	"os"
	"strings"
	"testing"
)

type stringWriter struct {
	strings.Builder
}

func (w *stringWriter) Write(p []byte) (int, error) {
	return w.Builder.Write(p)
}

func TestParseArgsUsage(t *testing.T) {
	// Test usage message when no args or -h is passed
	argsList := [][]string{
		{},
		{"-h"},
	}

	for _, args := range argsList {
		var buf strings.Builder
		exitCode := parseArgs(args, &buf)
		if exitCode != 1 {
			t.Errorf("Expected exit code 1 for args %v, got %d", args, exitCode)
		}
		output := buf.String()
		if output == "" || !strings.Contains(output, "Usage: mtojl markdown-file") {
			t.Errorf("Usage message missing for args %v", args)
		}
	}
}

func TestCheckFileExists(t *testing.T) {
	// Create a temp file to test existence
	tmpFile, err := os.CreateTemp("", "testfile.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if !checkFileExists(tmpFile.Name()) {
		t.Errorf("Expected file %s to exist", tmpFile.Name())
	}

	if checkFileExists("nonexistentfile.md") {
		t.Errorf("Expected nonexistent file to return false")
	}
}
