package parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ExtractAll extracts all files to the specified directory
func (p *ArchiveParser) ExtractAll(outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	if err := os.MkdirAll(outputDir, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
	}

	for i, entry := range p.entries {
		if err := p.extractEntry(entry, i, outputDir); err != nil {
			log.Printf("Warning: failed to extract entry %d: %v", i, err)
			continue
		}
	}

	return nil
}

// ExtractEntry extracts a single entry to the specified directory
func (p *ArchiveParser) ExtractEntry(entry FileEntry, outputDir string) error {
	if err := os.MkdirAll(outputDir, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
	}

	return p.extractEntry(entry, 0, outputDir)
}

// extractEntry extracts a single entry
func (p *ArchiveParser) extractEntry(entry FileEntry, index int, outputDir string) error {
	filename := p.generateFilename(entry, index)
	outputPath := filepath.Join(outputDir, filename)

	if len(entry.Content) == 0 {
		fmt.Printf("Skipping empty file: %s\n", filename)
		return nil
	}

	if err := os.WriteFile(outputPath, entry.Content, DefaultFilePerm); err != nil {
		return fmt.Errorf("failed to write file %q: %w", filename, err)
	}

	fmt.Printf("Extracted: %s (%d bytes)\n", filename, len(entry.Content))

	// Verify SHA1 if provided
	if entry.SHA1 != "" {
		if err := p.verifySHA1(entry.Content, entry.SHA1); err != nil {
			log.Printf("Warning: SHA1 verification failed for %s: %v", filename, err)
		} else {
			fmt.Printf("  ✓ SHA1 verified\n")
		}
	}

	return nil
}

// generateFilename generates a safe filename for the entry
func (p *ArchiveParser) generateFilename(entry FileEntry, index int) string {
	filename := entry.Filename
	if filename == "" {
		filename = fmt.Sprintf("file_%d_%s%s", index, entry.DocType, entry.Extension)
	}

	// Clean filename for filesystem
	replacements := map[string]string{
		"/": "_", "\\": "_", ":": "_", "*": "_", "?": "_",
		"\"": "_", "<": "_", ">": "_", "|": "_",
	}

	for old, new := range replacements {
		filename = strings.ReplaceAll(filename, old, new)
	}

	return filename
}
