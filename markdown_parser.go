package main

import (
	"bufio"
	"os"
	"regexp"
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

	// Regex to detect numbered list items like "1. ", "2. ", etc.
	numberedListRegex := regexp.MustCompile(`^\d+\.\s+`)

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

		// Detect list items (simple - or * with optional indentation) and numbered lists
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || numberedListRegex.MatchString(trimmed) {
			// Count leading spaces for depth, considering tabs as well
			leadingSpaces := 0
			for _, ch := range line {
				if ch == ' ' {
					leadingSpaces++
				} else if ch == '\t' {
					leadingSpaces += 2 // assuming tab width of 2 spaces
				} else {
					break
				}
			}

			depth := leadingSpaces/2 + 1
			content := trimmed
			if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
				content = strings.TrimSpace(trimmed[2:])
			} else if numberedListRegex.MatchString(trimmed) {
				idx := strings.Index(trimmed, ".")
				if idx != -1 && len(trimmed) > idx+1 {
					content = strings.TrimSpace(trimmed[idx+1:])
				}
			}

			// Add depth info to list item path by appending list depth
			listPath := append([]string{}, headingStack...)
			// Clear previous list-item entries if any
			for len(listPath) > 0 && listPath[len(listPath)-1] == "list-item" {
				listPath = listPath[:len(listPath)-1]
			}
			// Append "list-item" depth times
			for i := 0; i < depth; i++ {
				listPath = append(listPath, "list-item")
			}
			elements = append(elements, &MarkdownElement{
				Type:    "list",
				Content: content,
				Depth:   depth,
				Path:    listPath,
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
