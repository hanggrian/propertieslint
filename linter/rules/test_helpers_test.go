package rules_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hanggrian/propertieslint/linter"
)

type Config = linter.Config

func DefaultConfig() Config {
	return linter.DefaultConfig()
}

func lintSample(t *testing.T, sample string) []linter.Issue {
	t.Helper()
	return lintSampleWithConfig(t, sample, DefaultConfig())
}

func lintSampleWithConfig(t *testing.T, sample string, config Config) []linter.Issue {
	t.Helper()
	issues, err := linter.LintReader("sample.properties", strings.NewReader(sample), config)
	if err != nil {
		t.Fatalf("LintReader returned error: %v", err)
	}
	return issues
}

func lintFileSample(t *testing.T, content string, config Config) []linter.Issue {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sample.properties")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	issues, err := linter.LintFile(path, config)
	if err != nil {
		t.Fatalf("LintFile returned error: %v", err)
	}
	return issues
}

func assertIssueAt(
	t *testing.T,
	issues []linter.Issue,
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

func assertNoIssueContaining(t *testing.T, issues []linter.Issue, messageSubstring string) {
	t.Helper()
	for _, issue := range issues {
		if strings.Contains(strings.ToLower(issue.Message), strings.ToLower(messageSubstring)) {
			t.Fatalf("did not expect issue containing %q, got %#v", messageSubstring, issues)
		}
	}
}
