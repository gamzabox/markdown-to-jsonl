package main

import "strings"

// optimizeCodeBlock performs simple normalization on the content of a code block:
// - trims leading/trailing empty lines
// - collapses multiple consecutive empty lines into a single empty line
// - trims trailing spaces on each line
func optimizeCodeBlock(block *MarkdownElement) *MarkdownElement {
	if block == nil {
		return block
	}
	lines := strings.Split(block.Content, "\n")

	// Trim leading empty lines
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	// Trim trailing empty lines
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	trimmed := lines[start:end]

	// Collapse consecutive empty lines and trim trailing spaces per line
	var out []string
	prevEmpty := false
	for _, ln := range trimmed {
		t := strings.TrimRight(ln, " \t")
		isEmpty := strings.TrimSpace(t) == ""
		if isEmpty {
			if prevEmpty {
				// skip additional empty line
				continue
			}
			out = append(out, "")
			prevEmpty = true
		} else {
			out = append(out, t)
			prevEmpty = false
		}
	}

	// Ensure content isn't empty; if so, keep original content (preserve)
	if len(out) == 0 {
		return block
	}

	block.Content = strings.Join(out, "\n") + "\n"
	return block
}
