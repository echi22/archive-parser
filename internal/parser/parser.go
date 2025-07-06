package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// New creates a new ArchiveParser instance
func New() *ArchiveParser {
	return &ArchiveParser{
		entries: make([]FileEntry, 0),
	}
}

// ParseFile parses an archive file from disk
func (p *ArchiveParser) ParseFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filename, err)
	}

	return p.ParseData(data)
}

// ParseData parses archive data from a byte slice
func (p *ArchiveParser) ParseData(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	p.data = data
	return p.parseContent()
}

// parseContent parses the main content into sections
func (p *ArchiveParser) parseContent() error {
	sections := strings.Split(string(p.data), SectionDelimiter)
	fmt.Printf("Found %d sections\n", len(sections))

	for i, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		sectionData := p.normalizeSectionData(section, i, len(sections))
		trimmed := strings.TrimSpace(sectionData)

		if p.shouldSkipSection(trimmed, i) {
			continue
		}

		entry, err := p.parseSection(sectionData)
		if err != nil {
			log.Printf("Warning: failed to parse section %d: %v", i, err)
			continue
		}

		if entry != nil {
			p.entries = append(p.entries, *entry)
		}
	}

	return nil
}

// normalizeSectionData normalizes section data based on position
func (p *ArchiveParser) normalizeSectionData(section string, index, totalSections int) string {
	switch {
	case index == totalSections-1:
		return strings.TrimSuffix(section, "**")
	default:
		return section
	}
}

// shouldSkipSection determines if a section should be skipped
func (p *ArchiveParser) shouldSkipSection(trimmed string, index int) bool {
	if !strings.HasPrefix(trimmed, DocumentPrefix) {
		fmt.Printf("Skipping non-document section %d (not %s): %s\n",
			index, DocumentPrefix, p.truncateString(trimmed, 50))
		return true
	}

	return false
}

// parseSection parses a single section into a FileEntry
func (p *ArchiveParser) parseSection(sectionData string) (*FileEntry, error) {
	entry := &FileEntry{
		Metadata: make(map[string]string),
	}

	sigIndex := strings.Index(sectionData, SignatureMarker)
	if sigIndex == -1 {
		return nil, fmt.Errorf("no signature marker found")
	}

	// Parse header
	headerData := sectionData[:sigIndex]
	p.parseHeader(headerData, entry)

	// Extract content - no cleaning, just use the length header
	content, err := p.extractSectionContent([]byte(sectionData), sigIndex, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content: %w", err)
	}

	entry.Content = content
	fmt.Printf("Extracted %s: %d bytes (declared length: %d)\n",
		entry.Filename, len(entry.Content), entry.ContentLengthHint)

	return entry, nil
}

// extractSectionContent extracts content from a section using length header
func (p *ArchiveParser) extractSectionContent(sectionBytes []byte, sigIndex int, entry *FileEntry) ([]byte, error) {
	sigOffset := sigIndex + len(SignatureMarker)

	if sigOffset+4 > len(sectionBytes) {
		return nil, fmt.Errorf("not enough bytes after signature marker for length prefix")
	}

	lengthBytes := sectionBytes[sigOffset : sigOffset+4]
	contentLength := binary.LittleEndian.Uint32(lengthBytes)
	entry.ContentLengthHint = contentLength
	entry.Metadata["ContentLengthHint"] = fmt.Sprintf("%d", contentLength)

	rawContent := sectionBytes[sigOffset+4:]

	// Extract exactly the number of bytes specified by the length header
	if int(contentLength) <= len(rawContent) {
		return rawContent[:contentLength], nil
	} else {
		log.Printf("Warning: declared content length (%d) exceeds available data (%d)",
			contentLength, len(rawContent))
		return rawContent, nil // Return what we have
	}
}

// GetEntries returns all parsed file entries
func (p *ArchiveParser) GetEntries() []FileEntry {
	return p.entries
}

// GetEntry returns a specific entry by index
func (p *ArchiveParser) GetEntry(index int) (*FileEntry, error) {
	if index < 0 || index >= len(p.entries) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(p.entries))
	}
	return &p.entries[index], nil
}

// GetEntryByFilename returns the first entry with the specified filename
func (p *ArchiveParser) GetEntryByFilename(filename string) (*FileEntry, error) {
	for i := range p.entries {
		if p.entries[i].Filename == filename {
			return &p.entries[i], nil
		}
	}
	return nil, fmt.Errorf("entry with filename %q not found", filename)
}

// Count returns the number of entries
func (p *ArchiveParser) Count() int {
	return len(p.entries)
}

// truncateString truncates a string to the specified length
func (p *ArchiveParser) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
