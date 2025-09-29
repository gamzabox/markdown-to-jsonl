package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// Import local packages to resolve undefined functions
)

// parseArgs parses CLI arguments and writes usage message to the provided writer if needed.
// Returns exit code: 0 if args are valid, 1 if usage message was printed.
func parseArgs(args []string, usageWriter *strings.Builder) int {
	if len(args) == 0 || (len(args) == 1 && args[0] == "-h") {
		usageWriter.WriteString("Markdown to JSONL 1.0.0\n")
		usageWriter.WriteString("Usage: mtojl markdown-file\n")
		return 1
	}
	return 0
}

// checkFileExists returns true if the given file path exists and is a regular file.
func checkFileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// getOutputFileName returns an output filename based on the input markdown filename.
// If a file with the base name exists, it appends a numeric postfix before the extension.

func getOutputFileName(inputFile string) string {
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	ext := ".jsonl"
	output := base + ext
	count := 0
	for {
		if !checkFileExists(output) {
			return output
		}
		count++
		output = fmt.Sprintf("%s-%d%s", base, count, ext)
	}
}

func main() {
	args := os.Args[1:]
	var usageBuilder strings.Builder
	exitCode := parseArgs(args, &usageBuilder)
	if exitCode == 1 {
		fmt.Print(usageBuilder.String())
		os.Exit(1)
	}

	inputFile := args[0]
	if !checkFileExists(inputFile) {
		fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	elements, err := parseMarkdownFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown file: %v\n", err)
		os.Exit(1)
	}

	for i, el := range elements {
		if el.Type == "codeblock" {
			elements[i] = optimizeCodeBlock(el)
		}
	}

	elements = doubleIndexJSONL(elements)

	outputFile := getOutputFileName(inputFile)
	err = writeJSONLOutput(outputFile, elements)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSONL output: %v\n", err)
		os.Exit(1)
	}
}
