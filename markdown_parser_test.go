package main

import (
	"os"
	"testing"
)

func TestParseMarkdownFile(t *testing.T) {
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

	elements, err := parseMarkdownFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("parseMarkdownFile error: %v", err)
	}

	// Basic checks on parsed elements
	if len(elements) == 0 {
		t.Fatal("Expected elements, got none")
	}

	// Check first element is heading 1
	if elements[0].Type != "heading" || elements[0].Content != "Heading 1" || elements[0].Depth != 1 {
		t.Errorf("First element incorrect: %+v", elements[0])
	}

	// Check code block content
	foundCodeBlock := false
	for _, el := range elements {
		if el.Type == "codeblock" {
			foundCodeBlock = true
			if el.Content == "" {
				t.Error("Code block content empty")
			}
		}
	}
	if !foundCodeBlock {
		t.Error("Expected to find a codeblock element")
	}
}
