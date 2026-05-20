package propertieslint

import (
	"strings"
	"testing"
)

func TestLintReaderReportsCommonIssues(t *testing.T) {
	const sample = "good=value\n" +
		"good=value\n" +
		"bad\\u12=value\n" +
		"duplicate=one\n" +
		"duplicate=two\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	if len(issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d: %#v", len(issues), issues)
	}

	assertIssueAt(t, issues, 3, 4, "invalid key escape")
	assertIssueAt(t, issues, 5, 1, "duplicate key")
}

func TestLintReaderFlagsMissingSeparatorAndValue(t *testing.T) {
	const sample = "keyOnly\n" +
		"colon:value\n" +
		"space value\n" +
		"empty=\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 1, 1, "missing key/value separator")
	assertIssueAt(t, issues, 2, 1, "missing key/value separator")
	assertIssueAt(t, issues, 3, 1, "missing key/value separator")
	assertIssueAt(t, issues, 4, 7, "missing value")
}

func TestLintFileFlagsUnterminatedContinuationSample(t *testing.T) {
	const sample = "foo=bar\\"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 1, 1, "unterminated line continuation")
}

func TestLintReaderCanDisableDuplicateKeyRule(t *testing.T) {
	config := DefaultConfig()
	config.DuplicateKey = false

	issues, err := lintReader("sample.properties", strings.NewReader("duplicate=one\nduplicate=two\n"), config)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	for _, issue := range issues {
		if strings.Contains(issue.Message, "duplicate key") {
			t.Fatalf("did not expect duplicate key issue, got %#v", issues)
		}
	}
}

func TestLintReaderReportsDuplicateKeyColumnAtStartOfKey(t *testing.T) {
	const sample = "duplicate=one\n duplicate=two\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 2, 2, "duplicate key")
}

func assertIssue(t *testing.T, issues []Issue, line int, messageSubstring string) {
	t.Helper()

	for _, issue := range issues {
		if issue.Line == line && strings.Contains(issue.Message, messageSubstring) {
			return
		}
	}

	t.Fatalf("expected issue on line %d containing %q, got %#v", line, messageSubstring, issues)
}

func assertIssueAt(t *testing.T, issues []Issue, line int, column int, messageSubstring string) {
	t.Helper()

	for _, issue := range issues {
		if issue.Line == line && issue.Column == column && strings.Contains(issue.Message, messageSubstring) {
			return
		}
	}

	t.Fatalf("expected issue on line %d column %d containing %q, got %#v", line, column, messageSubstring, issues)
}

func TestLintReaderFlagsUntrimmedEntry(t *testing.T) {
	const sample = "key=value\n" +
		" key=value\n" +
		"key =value\n" +
		"key= value\n" +
		"key=value \n"

	config := DefaultConfig()
	config.DuplicateKey = false

	issues, err := lintReader("sample.properties", strings.NewReader(sample), config)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 2, 1, "key has leading whitespace")
	assertIssueAt(t, issues, 3, 4, "key has trailing whitespace")
	assertIssueAt(t, issues, 4, 5, "value has leading whitespace")
	assertIssueAt(t, issues, 5, 10, "value has trailing whitespace")
}

func TestLintReaderFlagsDuplicateBlankLine(t *testing.T) {
	const sample = "key=value\n" +
		"\n" +
		"\n" +
		"key2=value2\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 3, 1, "duplicate blank line")
}

func TestLintReaderFlagsNoLeadingBlankLine(t *testing.T) {
	const sample = "\n" +
		"key=value\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssueAt(t, issues, 1, 1, "leading blank line")
}

func TestLintFileFlagsTrailingNewline(t *testing.T) {
	const sample = "key=value"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	if len(issues) > 0 {
		t.Fatalf("lintReader should not report issues for missing trailing newline, got %#v", issues)
	}
}

func TestLintReaderCanDisableUntrimmedEntry(t *testing.T) {
	const sample = " key = value \n"

	config := DefaultConfig()
	config.UntrimmedEntry = false

	issues, err := lintReader("sample.properties", strings.NewReader(sample), config)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	for _, issue := range issues {
		if strings.Contains(issue.Message, "whitespace") {
			t.Fatalf("did not expect whitespace issue, got %#v", issues)
		}
	}
}

func TestLintReaderCanDisableDuplicateBlankLine(t *testing.T) {
	const sample = "key=value\n\n\n"

	config := DefaultConfig()
	config.DuplicateBlankLine = false

	issues, err := lintReader("sample.properties", strings.NewReader(sample), config)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	for _, issue := range issues {
		if strings.Contains(issue.Message, "duplicate blank line") {
			t.Fatalf("did not expect duplicate blank line issue, got %#v", issues)
		}
	}
}
