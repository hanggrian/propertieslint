package linter

import "strings"

func untrimmedEntryIssue(path string, line int, key string, value string) []*Issue {
	var issues []*Issue

	if len(key) > 0 &&
		(key[0] == ' ' || key[0] == '\t' || key[0] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  firstWhitespaceColumn(key),
			Message: "key has leading whitespace",
		})
	}

	if len(key) > 0 &&
		(key[len(key)-1] == ' ' || key[len(key)-1] == '\t' || key[len(key)-1] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  lastNonWhitespaceColumn(key) + 1,
			Message: "key has trailing whitespace",
		})
	}

	if len(value) > 0 &&
		(value[0] == ' ' || value[0] == '\t' || value[0] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  len(key) + 2,
			Message: "value has leading whitespace",
		})
	}

	if len(value) > 0 &&
		(value[len(value)-1] == ' ' || value[len(value)-1] == '\t' || value[len(value)-1] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  len(key) + 2 + len(strings.TrimRight(value, " \t\f")),
			Message: "value has trailing whitespace",
		})
	}

	return issues
}
