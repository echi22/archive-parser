package parser

import (
	"strings"
	"unicode"
)

// parseHeader parses the header section for metadata
func (p *ArchiveParser) parseHeader(headerData string, entry *FileEntry) {
	lines := strings.Split(headerData, "\n")

	// Parse first line for document type
	if len(lines) > 0 {
		p.parseDocumentType(lines[0], entry)
	}

	// Parse metadata lines
	for _, line := range lines {
		p.parseMetadataLine(line, entry)
	}
}

// parseDocumentType extracts document type from the first line
func (p *ArchiveParser) parseDocumentType(firstLine string, entry *FileEntry) {
	trimmed := strings.TrimSpace(firstLine)
	if strings.HasPrefix(trimmed, DocumentPrefix) {
		docType := strings.TrimPrefix(trimmed, DocumentPrefix)
		docType = strings.TrimFunc(docType, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		if docType != "" {
			entry.DocType = docType
		}
	}
}

// parseMetadataLine parses a single metadata line
func (p *ArchiveParser) parseMetadataLine(line string, entry *FileEntry) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "%%") {
		return
	}

	// Parse key/value pairs
	if strings.Contains(line, "/") && !strings.HasPrefix(line, SignatureMarker) {
		parts := strings.SplitN(line, "/", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			entry.Metadata[key] = value
			p.mapMetadataToFields(key, value, entry)
		}
	}
}

// mapMetadataToFields maps metadata keys to specific entry fields
func (p *ArchiveParser) mapMetadataToFields(key, value string, entry *FileEntry) {
	switch key {
	case "ENV_GUID":
		entry.EnvGUID = value
	case "EXT":
		entry.Extension = value
	case "FILENAME":
		entry.Filename = value
	case "GUID":
		entry.GUID = value
	case "SHA1":
		entry.SHA1 = value
	case "TYPE":
		entry.Type = value
	case "DOCTYPE":
		if entry.DocType == "" {
			entry.DocType = value
		}
	}
}
