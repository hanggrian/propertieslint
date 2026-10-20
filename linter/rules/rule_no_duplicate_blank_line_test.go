package rules_test

import "testing"

func TestNoDuplicateBlankLineFlagsSecondConsecutiveLine(t *testing.T) {
	issues := lintSample(t, "key=value\n\n\nkey2=value2\n")

	assertIssueAt(t, issues, 3, 1, "Remove duplicate blank line.")
}

func TestNoDuplicateBlankLineFlagsEveryExtraLine(t *testing.T) {
	issues := lintSample(t, "key=value\n\n\n\nkey2=value2\n")

	assertIssueAt(t, issues, 3, 1, "Remove duplicate blank line.")
	assertIssueAt(t, issues, 4, 1, "Remove duplicate blank line.")
}

func TestNoDuplicateBlankLineFlagsLeadingDuplicates(t *testing.T) {
	issues := lintSample(t, "\n\nkey=value\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate blank line.")
}

func TestNoDuplicateBlankLineTreatsWhitespaceOnlyLinesAsBlank(t *testing.T) {
	issues := lintSample(t, "key=value\n   \n\t\nkey2=value2\n")

	assertIssueAt(t, issues, 3, 1, "Remove duplicate blank line.")
}

func TestNoDuplicateBlankLineAllowsSeparatedBlankLines(t *testing.T) {
	issues := lintSample(t, "key=value\n\n# comment\n\nkey2=value2\n")

	assertNoIssueContaining(t, issues, "duplicate blank line")
}

func TestNoDuplicateBlankLineAllowsSingleBlankLine(t *testing.T) {
	issues := lintSample(t, "key=value\n\nkey2=value2\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoDuplicateBlankLineCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoDuplicateBlankLine = false

	issues := lintSampleWithConfig(t, "key=value\n\n\n", config)

	assertNoIssueContaining(t, issues, "duplicate blank line")
}
