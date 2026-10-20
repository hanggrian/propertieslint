package linter

import (
	"strings"
	"testing"
)

func assertIssueAt(
	t *testing.T,
	issues []Issue,
	line int,
	column int,
	messageSubstring string,
) {
	t.Helper()
	for _, issue := range issues {
		if issue.Line == line &&
			issue.Column == column &&
			strings.Contains(strings.ToLower(issue.Message), strings.ToLower(messageSubstring)) {
			return
		}
	}
	t.Fatalf(
		"expected issue on line %d column %d containing %q, got %#v",
		line,
		column,
		messageSubstring,
		issues,
	)
}
