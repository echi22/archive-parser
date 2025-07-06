package parser

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeader(t *testing.T) {
	p := New()
	entry := &FileEntry{Metadata: map[string]string{}}
	header := "DOCUTestType\nFILENAME/test.txt\nGUID/1234\nSHA1/abcd"
	p.parseHeader(header, entry)

	assert.Equal(t, "TestType", entry.DocType)
	assert.Equal(t, "test.txt", entry.Filename)
	assert.Equal(t, "1234", entry.GUID)
	assert.Equal(t, "abcd", entry.SHA1)
	assert.Equal(t, "test.txt", entry.Metadata["FILENAME"])
}

func TestParseData(t *testing.T) {
	p := New()
	// Compose raw data: DOCU section + signature + content length + content
	raw := []byte("DOCUTest\nFILENAME/test.txt\n_SIG/D.C." + string([]byte{2, 0, 0, 0}) + "Hi**")
	err := p.ParseData(raw)
	assert.NoError(t, err)
	assert.Equal(t, 1, p.Count())

	entry := p.GetEntries()[0]
	assert.Equal(t, "test.txt", entry.Filename)
	assert.Equal(t, []byte("Hi"), entry.Content)
}

func TestVerifySHA1(t *testing.T) {
	p := New()
	data := []byte("hello")
	hash := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" // sha1("hello")

	err := p.verifySHA1(data, hash)
	if err != nil {
		t.Errorf("SHA1 should verify, but got error: %v", err)
	}
}

func TestExtractEntry(t *testing.T) {
	dir := t.TempDir()
	p := New()
	entry := FileEntry{Filename: "test.txt", Content: []byte("hi")}

	err := p.ExtractEntry(entry, dir)
	if err != nil {
		t.Fatalf("ExtractEntry failed: %v", err)
	}

	path := filepath.Join(dir, "test.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(data) != "hi" {
		t.Errorf("expected 'hi', got '%s'", string(data))
	}
}

func TestExtractSectionContent(t *testing.T) {
	p := New()
	entry := &FileEntry{Metadata: map[string]string{}}
	section := "Header\n_SIG/D.C."
	buf := []byte(section)
	buf = append(buf, make([]byte, 4)...) // 4-byte length
	binary.LittleEndian.PutUint32(buf[len(buf)-4:], 5)
	buf = append(buf, []byte("abcde")...)

	content, err := p.extractSectionContent(buf, strings.Index(section, SignatureMarker), entry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(content) != "abcde" {
		t.Errorf("expected content 'abcde', got '%s'", string(content))
	}
}

func TestGenerateFilename(t *testing.T) {
	p := New()
	entry := FileEntry{DocType: "X", Extension: ".txt"}
	name := p.generateFilename(entry, 0)
	if !strings.HasSuffix(name, ".txt") {
		t.Errorf("filename does not have extension: %s", name)
	}
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		t.Errorf("filename contains unsafe characters: %s", name)
	}
}

func TestParseData_EmptyInput(t *testing.T) {
	p := New()
	err := p.ParseData([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be empty")
}

func TestParseData_MissingSignature(t *testing.T) {
	p := New()
	raw := []byte("DOCUTest\nFILENAME/test.txt\n")
	err := p.ParseData(raw)
	assert.NoError(t, err)
	assert.Equal(t, 0, p.Count()) // No valid entry due to missing _SIG
}

func TestParseData_InvalidLengthBytes(t *testing.T) {
	p := New()
	raw := []byte("DOCUTest\nFILENAME/test.txt\n_SIG/D.C.") // no 4 bytes after _SIG
	err := p.ParseData(raw)
	assert.NoError(t, err)
	assert.Equal(t, 0, p.Count())
}

func TestParseData_TooLongDeclaredLength(t *testing.T) {
	p := New()
	raw := append([]byte("DOCUTest\nFILENAME/test.txt\n_SIG/D.C."), []byte{10, 0, 0, 0}...)
	raw = append(raw, []byte("short")...)
	raw = append(raw, []byte("**")...)

	err := p.ParseData(raw)
	assert.NoError(t, err)
	assert.Equal(t, 1, p.Count())
	entry := p.GetEntries()[0]
	assert.Equal(t, []byte("short"), entry.Content)
}

func TestExtractAllAndEntry(t *testing.T) {
	p := New()
	entry := FileEntry{
		Filename: "hello.txt",
		Content:  []byte("hello world"),
		SHA1:     "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", // valid SHA1
		Metadata: map[string]string{},
	}
	p.entries = append(p.entries, entry)

	tempDir := t.TempDir()
	err := p.ExtractAll(tempDir)
	assert.NoError(t, err)

	// Check file contents
	data, err := os.ReadFile(tempDir + "/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), data)
}

func TestExtractEntry_EmptyContent(t *testing.T) {
	p := New()
	entry := FileEntry{
		Filename: "empty.txt",
		Content:  []byte{},
		Metadata: map[string]string{},
	}

	tempDir := t.TempDir()
	err := p.ExtractEntry(entry, tempDir)
	assert.NoError(t, err) // Should skip without error

	_, err = os.Stat(tempDir + "/empty.txt")
	assert.Error(t, err) // File should not exist
}

func TestParseData_MultipleSections(t *testing.T) {
	p := New()

	// Compose raw data with two sections separated by SectionDelimiter "**%%"
	section1 := []byte("DOCUType1\nFILENAME/file1.txt\n_SIG/D.C.")
	section1 = append(section1, []byte{5, 0, 0, 0}...) // content length 5
	section1 = append(section1, []byte("Hello")...)

	section2 := []byte("DOCUType2\nFILENAME/file2.txt\n_SIG/D.C.")
	section2 = append(section2, []byte{3, 0, 0, 0}...) // content length 3
	section2 = append(section2, []byte("Bye")...)

	// Join sections with the SectionDelimiter
	raw := append(section1, []byte("**%%")...)
	raw = append(raw, section2...)
	raw = append(raw, []byte("**")...) // end marker for last section

	err := p.ParseData(raw)
	assert.NoError(t, err)
	assert.Equal(t, 2, p.Count())

	entry1 := p.GetEntries()[0]
	assert.Equal(t, "file1.txt", entry1.Filename)
	assert.Equal(t, []byte("Hello"), entry1.Content)

	entry2 := p.GetEntries()[1]
	assert.Equal(t, "file2.txt", entry2.Filename)
	assert.Equal(t, []byte("Bye"), entry2.Content)
}

