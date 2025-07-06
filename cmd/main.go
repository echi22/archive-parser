package main

import (
	"fmt"
	"log"
	"os"

	"github.com/archive-parser/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <archive-file> [output-directory]")
		os.Exit(1)
	}

	archiveFile := os.Args[1]
	outputDir := "extracted"
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	archiveParser := parser.New()

	fmt.Printf("Parsing archive: %s\n", archiveFile)
	if err := archiveParser.ParseFile(archiveFile); err != nil {
		log.Fatalf("Failed to parse archive: %v", err)
	}

	archiveParser.PrintSummary()

	fmt.Printf("\nExtracting files to: %s\n", outputDir)
	if err := archiveParser.ExtractAll(outputDir); err != nil {
		log.Fatalf("Failed to extract files: %v", err)
	}

	fmt.Println("\nExtraction completed!")
}
