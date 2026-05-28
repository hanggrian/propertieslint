package linter

import (
	"fmt"
	"strings"
)

func missingSeparatorIssue(
	path string,
	line int,
	column int,
	separatorFound bool,
) *Issue {
	if separatorFound {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Missing key/value separator.",
	}
}

func untrimmedEntryIssue(
	path string,
	line int,
	key string,
	value string,
) []*Issue {
	var issues []*Issue
	if len(key) > 0 &&
		(key[0] == ' ' || key[0] == '\t' || key[0] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  firstWhitespaceColumn(key),
			Message: "Key has leading whitespace.",
		})
	}
	if len(key) > 0 &&
		(key[len(key)-1] == ' ' || key[len(key)-1] == '\t' || key[len(key)-1] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  lastNonWhitespaceColumn(key) + 1,
			Message: "Key has trailing whitespace.",
		})
	}
	if len(value) > 0 &&
		(value[0] == ' ' || value[0] == '\t' || value[0] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  len(key) + 2,
			Message: "Value has leading whitespace.",
		})
	}
	if len(value) > 0 &&
		(value[len(value)-1] == ' ' || value[len(value)-1] == '\t' || value[len(value)-1] == '\f') {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  len(key) + 2 + len(strings.TrimRight(value, " \t\f")),
			Message: "Value has trailing whitespace.",
		})
	}
	return issues
}

func validateEscapes(
	path string,
	line int,
	rawLogical string,
	key string,
	value string,
) (string, []Issue, bool) {
	rawKey, _, _ := splitKeyValueRaw(rawLogical)
	normalizedKey, err := unescape(key)
	if err != nil {
		return "", []Issue{{
			Path:    path,
			Line:    line,
			Column:  firstBackslashColumn(rawKey),
			Message: fmt.Sprintf("Invalid key escape: %v.", err),
		}}, false
	}
	if _, err := unescape(value); err != nil {
		return normalizedKey, []Issue{{
			Path:    path,
			Line:    line,
			Column:  len(rawKey) + 2 + firstBackslashColumn(value) - 1,
			Message: fmt.Sprintf("Invalid value escape: %v.", err),
		}}, true
	}
	return normalizedKey, nil, true
}
