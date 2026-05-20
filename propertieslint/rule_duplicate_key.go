package propertieslint

import "fmt"

func duplicateKeyIssue(path string, line int, key string, seenKeys map[string]int) *Issue {
	if previousLine, ok := seenKeys[key]; ok {
		return &Issue{
			Path:    path,
			Line:    line,
			Message: fmt.Sprintf("duplicate key (first seen at line %d)", previousLine),
		}
	}
	seenKeys[key] = line
	return nil
}
