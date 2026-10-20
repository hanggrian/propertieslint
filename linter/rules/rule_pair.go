package rules

import (
	"fmt"
	"strings"
	"unicode"
)

func KeyNameIssue(
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

func KeyOrderIssue(
	path string,
	line int,
	column int,
	previousKey string,
	key string,
) *Issue {
	if previousKey == "" ||
		strings.Compare(strings.ToLower(previousKey), strings.ToLower(key)) <= 0 {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: fmt.Sprintf("Arrange '%s' before '%s'", key, previousKey),
	}
}

func NoDuplicateKeyIssue(
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
			Message: fmt.Sprintf("Remove duplicate key (first seen at line %d).", previousLine),
		}
	}
	seenKeys[key] = line
	return nil
}

func NoEmptyKeyIssue(
	path string,
	line int,
) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  1,
		Message: "Add key.",
	}
}

func NoEmptyValueIssue(
	path string,
	line int,
	column int,
	separatorFound bool,
	value string,
) *Issue {
	if !separatorFound ||
		value != "" {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Add value.",
	}
}
