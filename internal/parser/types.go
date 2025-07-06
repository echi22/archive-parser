package parser

// FileEntry represents a file within the archive
type FileEntry struct {
	DocType           string
	EnvGUID           string
	Extension         string
	Filename          string
	GUID              string
	SHA1              string
	Type              string
	Content           []byte
	Metadata          map[string]string
	ContentLengthHint uint32
}

// ArchiveParser handles parsing of the custom .env format
type ArchiveParser struct {
	entries []FileEntry
	data    []byte
}
