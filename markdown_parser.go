package main

import (
	"bufio"
	"os"
	"strings"
)

// MarkdownElement represents a parsed element from markdown with hierarchy info.
type MarkdownElement struct {
	Type     string // e.g., "heading", "list", "codeblock", "text"
	Content  string
	Depth    int // hierarchy depth for headings/lists
	Children []*MarkdownElement
	Path     []string // path of headings leading to this element
}

// parseMarkdownFile parses the markdown file preserving heading/list hierarchy and code blocks.
func parseMarkdownFile(filePath string) ([]*MarkdownElement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var elements []*MarkdownElement
	var headingStack []string
	var currentCodeBlock *MarkdownElement

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Detect code block start/end
		if strings.HasPrefix(trimmed, "```") {
			if currentCodeBlock == nil {
				currentCodeBlock = &MarkdownElement{
					Type:    "codeblock",
					Content: "",
					Depth:   len(headingStack),
					Path:    append([]string{}, headingStack...),
				}
			} else {
				elements = append(elements, currentCodeBlock)
				currentCodeBlock = nil
			}
			continue
		}

		if currentCodeBlock != nil {
			currentCodeBlock.Content += line + "\n"
			continue
		}

		// Detect headings
		if strings.HasPrefix(trimmed, "#") {
			depth := 0
			for i := 0; i < len(trimmed) && trimmed[i] == '#'; i++ {
				depth++
			}
			content := strings.TrimSpace(trimmed[depth:])
			// Adjust heading stack
			if depth <= len(headingStack) {
				headingStack = headingStack[:depth-1]
			}
			headingStack = append(headingStack, content)
			elements = append(elements, &MarkdownElement{
				Type:    "heading",
				Content: content,
				Depth:   depth,
				Path:    append([]string{}, headingStack...),
			})
			continue
		}

		// Detect list items (simple - or * with optional indentation)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			// Count leading spaces for depth
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			depth := leadingSpaces/2 + 1
			content := strings.TrimSpace(trimmed[2:])
			elements = append(elements, &MarkdownElement{
				Type:    "list",
				Content: content,
				Depth:   depth,
				Path:    append([]string{}, headingStack...),
			})
			continue
		}

		// Plain text or other lines
		elements = append(elements, &MarkdownElement{
			Type:    "text",
			Content: trimmed,
			Depth:   len(headingStack),
			Path:    append([]string{}, headingStack...),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return elements, nil
}
