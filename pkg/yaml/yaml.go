package yaml

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// SplitDocuments returns a slice of all documents contained in a YAML string. Multiple documents can be divided by the
// YAML document separator (---). It allows for white space and comments to be after the separator on the same line,
// but will return an error if anything else is on the line.
func SplitDocuments(s string) ([]string, error) {
	docs := make([]string, 0)
	if len(s) > 0 {
		// The YAML document separator is any line that starts with ---
		yamlSeparatorRegexp := regexp.MustCompile(`\n---.*\n`)

		// Find all separators, check them for invalid content, and append each document to docs
		separatorLocations := yamlSeparatorRegexp.FindAllStringIndex(s, -1)
		prev := 0
		for i := range separatorLocations {
			loc := separatorLocations[i]
			separator := s[loc[0]:loc[1]]
			// If the next non-whitespace character on the line following the separator is not a comment, return an error
			trimmedContentAfterSeparator := strings.TrimSpace(separator[4:])
			if len(trimmedContentAfterSeparator) > 0 && trimmedContentAfterSeparator[0] != '#' {
				return nil, fmt.Errorf("invalid document separator: %s", strings.TrimSpace(separator))
			}
			// Remove all whitespace
			result := s[prev:loc[0]]
			if len(result) > 0 && !isAllWhitespace(result) {
				docs = append(docs, result)
			}
			prev = loc[1]
		}
		docs = append(docs, s[prev:])
	}
	return docs, nil
}

func isAllWhitespace(str string) bool {
	for _, r := range str {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
