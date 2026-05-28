package linter

import (
	"strings"
	"testing"
)

func TestLintReaderReportsCommonIssues(t *testing.T) {
	const sample = "good=value\n" +
		"other=value\n" +
		"UPPER=value\n" +
		"quoted=\"bar\"\n" +
		"bad\\u12=value\n" +
		"duplicate=one\n" +
		"duplicate=two\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	if len(issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d: %#v", len(issues), issues)
	}
	assertIssueAt(t, issues, 3, 1, "Key name cannot be all uppercase")
	assertIssueAt(t, issues, 4, 8, "Value should not be quoted")
	assertIssueAt(t, issues, 5, 4, "Invalid key escape")
	assertIssueAt(t, issues, 7, 1, "Duplicate key")
}

func TestLintReaderFlagsMissingSeparatorAndValue(t *testing.T) {
	const sample = "keyOnly\n" +
		"colon:value\n" +
		"space value\n" +
		"empty=\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 1, "Missing key/value separator")
	assertIssueAt(t, issues, 2, 1, "Missing key/value separator")
	assertIssueAt(t, issues, 3, 1, "Missing key/value separator")
	assertIssueAt(t, issues, 4, 7, "Missing value")
}

func TestLintFileFlagsUnterminatedContinuationSample(t *testing.T) {
	const sample = "foo=bar\\"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 1, "Unterminated line continuation")
}

func TestLintReaderCanDisableDuplicateKeyRule(t *testing.T) {
	config := DefaultConfig()
	config.DuplicateKey = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("duplicate=one\nduplicate=two\n"),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(strings.ToLower(issue.Message), "duplicate key") {
			t.Fatalf("did not expect duplicate key issue, got %#v", issues)
		}
	}
}

func TestLintReaderReportsDuplicateKeyColumnAtStartOfKey(t *testing.T) {
	const sample = "duplicate=one\n duplicate=two\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 2, 2, "Duplicate key")
}

func assertIssue(t *testing.T, issues []Issue, line int, messageSubstring string) {
	t.Helper()
	for _, issue := range issues {
		if issue.Line == line &&
			strings.Contains(strings.ToLower(issue.Message), strings.ToLower(messageSubstring)) {
			return
		}
	}
	t.Fatalf("expected issue on line %d containing %q, got %#v", line, messageSubstring, issues)
}

func assertIssueAt(t *testing.T, issues []Issue, line int, column int, messageSubstring string) {
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

func TestLintReaderFlagsUntrimmedEntry(t *testing.T) {
	const sample = "key=value\n" +
		" key=value\n" +
		"key =value\n" +
		"key= value\n" +
		"key=value \n"
	config := DefaultConfig()
	config.DuplicateKey = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 2, 1, "Key has leading whitespace")
	assertIssueAt(t, issues, 3, 4, "Key has trailing whitespace")
	assertIssueAt(t, issues, 4, 5, "Value has leading whitespace")
	assertIssueAt(t, issues, 5, 10, "Value has trailing whitespace")
}

func TestLintReaderFlagsUntrimmedKey(t *testing.T) {
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(" key=value\n"),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 1, "Key has leading whitespace")
}

func TestLintReaderFlagsUntrimmedValue(t *testing.T) {
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("key= value\n"),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 5, "Value has leading whitespace")
}

func TestLintReaderFlagsDuplicateBlankLine(t *testing.T) {
	const sample = "key=value\n" +
		"\n" +
		"\n" +
		"key2=value2\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 3, 1, "Duplicate blank line")
}

func TestLintReaderFlagsCommentStyle(t *testing.T) {
	const sample = "!notacomment\n#good=value\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 1, "Illegal comment '!'")
}

func TestLintReaderCanDisableCommentStyle(t *testing.T) {
	config := DefaultConfig()
	config.CommentStyle = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("!notacomment\n"),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(strings.ToLower(issue.Message), "Illegal comment '!'") {
			t.Fatalf("did not expect comment-style issue, got %#v", issues)
		}
	}
}

func TestLintReaderFlagsCommentSpaces(t *testing.T) {
	const sample = "key=value#bad\nkey=value # good\nkey=value  #bad\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 10, "One space around inline '#'.")
	assertIssueAt(t, issues, 3, 12, "One space around inline '#'.")
}

func TestLintReaderCanDisableCommentSpaces(t *testing.T) {
	config := DefaultConfig()
	config.CommentSpaces = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("key=value#bad\nkey=value  #bad\n"),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(
			strings.ToLower(issue.Message),
			"One space around inline '#'.",
		) ||
			strings.Contains(
				strings.ToLower(issue.Message),
				"No space before and one space after '#'.",
			) {
			t.Fatalf("did not expect comment-spaces issue, got %#v", issues)
		}
	}
}

func TestLintReaderFlagsNoLeadingBlankLine(t *testing.T) {
	const sample = "\n" +
		"key=value\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 1, "Unexpected leading blank line.")
}

func TestLintFileFlagsTrailingNewline(t *testing.T) {
	const sample = "key=value"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	if len(issues) > 0 {
		t.Fatalf(
			"lintReader should not report issues for missing trailing newline, got %#v",
			issues,
		)
	}
}

func TestLintReaderCanDisableUntrimmedEntry(t *testing.T) {
	const sample = " key = value \n"
	config := DefaultConfig()
	config.UntrimmedEntry = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			config,
		)
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
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(issue.Message, "duplicate blank line") {
			t.Fatalf("did not expect duplicate blank line issue, got %#v", issues)
		}
	}
}

func TestLintReaderFlagsQuotedValue(t *testing.T) {
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("foo=\"bar\"\n"),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	assertIssueAt(t, issues, 1, 5, "Value should not be quoted")
}

func TestLintReaderCanDisableQuotedValueRule(t *testing.T) {
	config := DefaultConfig()
	config.QuotedValue = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("foo=\"bar\"\n"),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(strings.ToLower(issue.Message), "quoted") {
			t.Fatalf("did not expect quoted value issue, got %#v", issues)
		}
	}
}

func TestLintReaderAllowsMixedAndSymbolHeavyKeyNames(t *testing.T) {
	const sample = "api-V2=value\n" +
		"123_456=value\n"
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader(sample),
			DefaultConfig(),
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(issue.Message, "all uppercase") {
			t.Fatalf("did not expect uppercase key issue, got %#v", issues)
		}
	}
}

func TestLintReaderCanDisableKeyNameRule(t *testing.T) {
	config := DefaultConfig()
	config.KeyName = false
	issues, err :=
		lintReader(
			"sample.properties",
			strings.NewReader("UPPER=value\n"),
			config,
		)
	if err != nil {
		t.Fatalf("lintReader returned error: %v", err)
	}
	for _, issue := range issues {
		if strings.Contains(issue.Message, "all uppercase") {
			t.Fatalf("did not expect uppercase key issue, got %#v", issues)
		}
	}
}
