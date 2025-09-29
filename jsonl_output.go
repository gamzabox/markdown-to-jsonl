package main

import (
	"encoding/json"
	"os"
)

// writeJSONLOutput writes the markdown elements as JSONL to the specified output file.
func writeJSONLOutput(outputFile string, elements []*MarkdownElement) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, el := range elements {
		// Skip empty text elements to ignore empty lines in markdown
		if el.Type == "text" && el.Content == "" {
			continue
		}
		if err := encoder.Encode(el); err != nil {
			return err
		}
	}
	return nil
}
