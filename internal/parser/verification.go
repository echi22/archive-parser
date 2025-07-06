package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// verifySHA1 verifies the SHA1 hash of content
func (p *ArchiveParser) verifySHA1(content []byte, expectedSHA1 string) error {
	if expectedSHA1 == "" {
		return fmt.Errorf("expected SHA1 cannot be empty")
	}

	hasher := sha1.New()
	hasher.Write(content)
	actualSHA1 := hex.EncodeToString(hasher.Sum(nil))

	if !strings.EqualFold(actualSHA1, expectedSHA1) {
		return fmt.Errorf("SHA1 mismatch: expected %s, got %s", expectedSHA1, actualSHA1)
	}

	return nil
}
