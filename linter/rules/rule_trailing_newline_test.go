package rules_test

import "testing"

func TestTrailingNewlineFlagsMissingNewline(t *testing.T) {
	issues := lintFileSample(t, "key=value", DefaultConfig())

	assertIssueAt(t, issues, 1, 1, "Add trailing newline")
}

func TestTrailingNewlineFlagsMissingNewlineOnLastLine(t *testing.T) {
	issues := lintFileSample(t, "key=value\nkey2=value2", DefaultConfig())

	assertIssueAt(t, issues, 2, 1, "Add trailing newline")
}

func TestTrailingNewlineAllowsNewline(t *testing.T) {
	issues := lintFileSample(t, "key=value\n", DefaultConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestTrailingNewlineAllowsBlankLineAtEnd(t *testing.T) {
	issues := lintFileSample(t, "key=value\n\n", DefaultConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestTrailingNewlineAllowsEmptyFile(t *testing.T) {
	issues := lintFileSample(t, "", DefaultConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestTrailingNewlineFlagsCarriageReturnOnly(t *testing.T) {
	issues := lintFileSample(t, "key=value\r", DefaultConfig())

	assertIssueAt(t, issues, 1, 1, "Add trailing newline")
}

func TestTrailingNewlineIsNotReportedByReader(t *testing.T) {
	issues := lintSample(t, "key=value")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestTrailingNewlineCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.TrailingNewline = false

	issues := lintFileSample(t, "key=value", config)

	assertNoIssueContaining(t, issues, "Add trailing newline")
}
