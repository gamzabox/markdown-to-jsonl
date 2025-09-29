package main

// doubleIndexJSONL performs simple double-indexing on JSONL entries extracted from markdown elements.
// Implementation:
// - Index: global sequential index for each emitted element (starting at 1)
// - SecondaryIndex: section-based index (increments on each heading, applied to elements under that heading)
// This is a lightweight implementation to support downstream grouping/lookup.
func doubleIndexJSONL(elements []*MarkdownElement) []*MarkdownElement {
	if elements == nil {
		return nil
	}
	idx := 1
	section := 0
	for _, el := range elements {
		// If this element is a heading, increment section counter and assign its own secondary index.
		if el.Type == "heading" {
			section++
			el.SecondaryIndex = section
		} else {
			// Non-heading elements inherit current section (0 if before first heading)
			el.SecondaryIndex = section
		}
		el.Index = idx
		idx++
	}
	return elements
}
