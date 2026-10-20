package linter

import (
	"strings"
	"testing"
)

func TestLintReaderReportsCommonIssues(t *testing.T) {
	const sample = "good=value\n" +
		"other=value\n" +
		"UPPER=value\n" +
		"bad\\u12=value\n" +
		"duplicate=one\n" +
		"duplicate=two\n"
	config := DefaultConfig()
	config.KeyOrder = false
	issues, err :=
		LintReader(
			"sample.properties",
			strings.NewReader(sample),
			config,
		)
	if err != nil {
		t.Fatalf("LintReader returned error: %v", err)
	}
	if len(issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d: %#v", len(issues), issues)
	}
	assertIssueAt(t, issues, 3, 1, "Key name cannot be all uppercase.")
	assertIssueAt(t, issues, 4, 4, "Invalid key escape 'short unicode escape'.")
	assertIssueAt(t, issues, 6, 1, "Remove duplicate key (first seen at line 5).")
}
