package main

import (
	"encoding/json"
	"os"
)

// writeJSONLOutput writes the markdown elements as a JSON array to the specified output file.
// The resulting file is a single .json containing a JSON array where each element corresponds
// to one MarkdownElement (previously emitted as JSONL).
func writeJSONLOutput(outputFile string, elements []*MarkdownElement) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Filter out empty text elements
	var out []*MarkdownElement
	for _, el := range elements {
		if el.Type == "text" && el.Content == "" {
			continue
		}
		out = append(out, el)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
