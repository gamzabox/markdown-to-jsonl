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
		if err := encoder.Encode(el); err != nil {
			return err
		}
	}
	return nil
}
