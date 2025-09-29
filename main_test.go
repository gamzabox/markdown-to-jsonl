package main

import (
	"os"
	"testing"
)

func TestFullConversionProcess(t *testing.T) {
	// Create a temporary markdown file for testing
	content := "# Heading 1\nSome text under heading 1.\n\n## Heading 2\n- List item 1\n- List item 2\n\n```\ncode block line 1\ncode block line 2\n```\n\nPlain text after code block.\n"
	tmpFile, err := os.CreateTemp("", "test.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Simulate full conversion process
	elements, err := parseMarkdownFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("parseMarkdownFile error: %v", err)
	}

	// Optimize code blocks
	for i, el := range elements {
		if el.Type == "codeblock" {
			elements[i] = optimizeCodeBlock(el)
		}
	}

	// Apply double indexing
	elements = doubleIndexJSONL(elements)

	// Generate output filename
	outputFile := getOutputFileName(tmpFile.Name())

	// Write JSONL output (placeholder, actual writing function to be implemented)
	err = writeJSONLOutput(outputFile, elements)
	if err != nil {
		t.Fatalf("writeJSONLOutput error: %v", err)
	}

	// Check output file exists
	if !checkFileExists(outputFile) {
		t.Errorf("Expected output file %s to exist", outputFile)
	}

	// Clean up output file
	os.Remove(outputFile)
}
