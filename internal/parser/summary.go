package parser

import "fmt"

// PrintSummary prints a summary of all entries in the archive
func (p *ArchiveParser) PrintSummary() {
	fmt.Printf("\n=== Archive Summary ===\n")
	fmt.Printf("Total entries: %d\n\n", len(p.entries))

	for i, entry := range p.entries {
		p.printEntryDetails(i+1, entry)
	}
}

// printEntryDetails prints details for a single entry
func (p *ArchiveParser) printEntryDetails(index int, entry FileEntry) {
	fmt.Printf("Entry %d:\n", index)
	fmt.Printf("  Document Type: %s\n", entry.DocType)
	fmt.Printf("  Filename: %s\n", entry.Filename)
	fmt.Printf("  Extension: %s\n", entry.Extension)
	fmt.Printf("  GUID: %s\n", entry.GUID)
	fmt.Printf("  Type: %s\n", entry.Type)
	fmt.Printf("  Content Size: %d bytes\n", len(entry.Content))

	if entry.SHA1 != "" {
		fmt.Printf("  SHA1: %s\n", entry.SHA1)
	}

	if len(entry.Metadata) > 0 {
		fmt.Printf("  Additional Metadata:\n")
		for key, value := range entry.Metadata {
			if !p.isStandardMetadataKey(key) {
				fmt.Printf("    %s: %s\n", key, value)
			}
		}
	}

	fmt.Println()
}

// isStandardMetadataKey checks if a metadata key is already displayed in standard fields
func (p *ArchiveParser) isStandardMetadataKey(key string) bool {
	standardKeys := map[string]bool{
		"ENV_GUID":          true,
		"EXT":               true,
		"FILENAME":          true,
		"GUID":              true,
		"SHA1":              true,
		"TYPE":              true,
		"DOCTYPE":           true,
		"ContentLengthHint": true,
	}
	return standardKeys[key]
}

// GetSummaryStats returns summary statistics about the archive
func (p *ArchiveParser) GetSummaryStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_entries"] = len(p.entries)

	// Count by extension
	extensions := make(map[string]int)
	docTypes := make(map[string]int)
	totalSize := 0

	for _, entry := range p.entries {
		if entry.Extension != "" {
			extensions[entry.Extension]++
		}
		if entry.DocType != "" {
			docTypes[entry.DocType]++
		}
		totalSize += len(entry.Content)
	}

	stats["extensions"] = extensions
	stats["document_types"] = docTypes
	stats["total_content_size"] = totalSize

	return stats
}
