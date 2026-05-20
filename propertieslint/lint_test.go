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

	assertIssue(t, issues, 3, "invalid key escape")
	assertIssue(t, issues, 5, "duplicate key")
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

	assertIssue(t, issues, 1, "missing key/value separator")
	assertIssue(t, issues, 2, "missing key/value separator")
	assertIssue(t, issues, 3, "missing key/value separator")
	assertIssue(t, issues, 4, "missing value")
}

func TestLintFileFlagsUnterminatedContinuationSample(t *testing.T) {
	issues, err := LintFile("../sample/unterminated_continuation.properties", DefaultConfig())
	if err != nil {
		t.Fatalf("LintFile returned error: %v", err)
	}

	assertIssue(t, issues, 1, "unterminated line continuation")
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

func assertIssue(t *testing.T, issues []Issue, line int, messageSubstring string) {
	t.Helper()

	for _, issue := range issues {
		if issue.Line == line && strings.Contains(issue.Message, messageSubstring) {
			return
		}
	}

	t.Fatalf("expected issue on line %d containing %q, got %#v", line, messageSubstring, issues)
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

	assertIssue(t, issues, 2, "key has leading whitespace")
	assertIssue(t, issues, 3, "key has trailing whitespace")
	assertIssue(t, issues, 4, "value has leading whitespace")
	assertIssue(t, issues, 5, "value has trailing whitespace")
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

	assertIssue(t, issues, 3, "duplicate blank line")
}

func TestLintReaderFlagsNoLeadingBlankLine(t *testing.T) {
	const sample = "\n" +
		"key=value\n"

	issues, err := lintReader("sample.properties", strings.NewReader(sample), DefaultConfig())
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}

	assertIssue(t, issues, 1, "leading blank line")
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
