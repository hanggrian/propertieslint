package linter

import (
	"fmt"
	"unicode"
)

func duplicateKeyIssue(
	path string,
	line int,
	column int,
	key string,
	seenKeys map[string]int,
) *Issue {
	if previousLine, ok := seenKeys[key]; ok {
		return &Issue{
			Path:    path,
			Line:    line,
			Column:  column,
			Message: fmt.Sprintf("Duplicate key (first seen at line %d).", previousLine),
		}
	}
	seenKeys[key] = line
	return nil
}

func keyNameIssue(
	path string,
	line int,
	column int,
	key string,
) *Issue {
	letters := 0
	for _, r := range key {
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		if unicode.IsLower(r) {
			return nil
		}
	}
	if letters == 0 {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Key name cannot be all uppercase.",
	}
}
