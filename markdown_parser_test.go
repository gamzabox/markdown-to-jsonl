package main

import (
	"os"
	"testing"
)

func testListDepth(t *testing.T, elements []*MarkdownElement) {
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
}

func TestParseMarkdownFile(t *testing.T) {
	// Create a temporary markdown file for testing
	backticks := "```"
	content := `# Heading 1
Some text under heading 1.

## Heading 2
- List item 1
  - Nested list item 1
- List item 2

1. Numbered list item 1
   1. Nested numbered list item 1
2. Numbered list item 2
` + backticks + `
code block line 1
code block line 2
` + backticks + `

Plain text after code block.
`
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

	// Run list depth tests
	testListDepth(t, elements)

	// Additional test: check depth relationship for numbered lists
	var numListItem1Depth, nestedNumListItem1Depth int
	for _, el := range elements {
		if el.Type == "list" {
			if el.Content == "Numbered list item 1" {
				numListItem1Depth = el.Depth
			}
			if el.Content == "Nested numbered list item 1" {
				nestedNumListItem1Depth = el.Depth
			}
		}
	}
	if numListItem1Depth == 0 || nestedNumListItem1Depth == 0 {
		t.Error("Could not find depths for numbered list items")
	} else if nestedNumListItem1Depth != numListItem1Depth+1 {
		t.Errorf("'Nested numbered list item 1' depth %d is not one greater than 'Numbered list item 1' depth %d", nestedNumListItem1Depth, numListItem1Depth)
	}
}
