package main

import (
	"os"
	"testing"
)

func TestParseMarkdownFile(t *testing.T) {
	// Create a temporary markdown file for testing
	content := "# Heading 1\nSome text under heading 1.\n\n## Heading 2\n- List item 1\n  - Nested list item 1\n- List item 2\n\n```\ncode block line 1\ncode block line 2\n```\n\nPlain text after code block.\n"
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

	// Check list item depth and path
	foundList := false
	var listItems []*MarkdownElement
	for _, el := range elements {
		if el.Type == "list" {
			foundList = true
			listItems = append(listItems, el)
			if el.Depth == 0 {
				t.Error("List item depth should be greater than 0")
			}
			if len(el.Path) == 0 {
				t.Error("List item path should not be empty")
			}
			// Additional check: depth should match the count of "list-item" in path suffix
			listItemCount := 0
			for _, p := range el.Path {
				if p == "list-item" {
					listItemCount++
				}
			}
			if listItemCount != el.Depth {
				t.Errorf("List item depth %d does not match count of 'list-item' in path %d", el.Depth, listItemCount)
			}
		}
	}
	if !foundList {
		t.Error("Expected to find list item elements")
	}

	// Additional test: check depth relationship between "List item 1" and "Nested list item 1"
	var listItem1Depth, nestedListItem1Depth int
	for _, el := range listItems {
		if el.Content == "List item 1" {
			listItem1Depth = el.Depth
		}
		if el.Content == "Nested list item 1" {
			nestedListItem1Depth = el.Depth
		}
	}
	if listItem1Depth == 0 || nestedListItem1Depth == 0 {
		t.Error("Could not find depths for 'List item 1' or 'Nested list item 1'")
	} else if nestedListItem1Depth != listItem1Depth+1 {
		t.Errorf("'Nested list item 1' depth %d is not one greater than 'List item 1' depth %d", nestedListItem1Depth, listItem1Depth)
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
