package rules_test

import "testing"

func TestNoLeadingBlankLineFlagsBlankFirstLine(t *testing.T) {
	issues := lintSample(t, "\nkey=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove leading blank line.")
}

func TestNoLeadingBlankLineFlagsWhitespaceOnlyFirstLine(t *testing.T) {
	issues := lintSample(t, "   \nkey=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove leading blank line.")
}

func TestNoLeadingBlankLineFlagsBeforeComment(t *testing.T) {
	issues := lintSample(t, "\n# comment\nkey=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove leading blank line.")
}

func TestNoLeadingBlankLineAllowsCommentFirst(t *testing.T) {
	issues := lintSample(t, "# comment\n\nkey=value\n")

	assertNoIssueContaining(t, issues, "leading blank line")
}

func TestNoLeadingBlankLineAllowsLaterBlankLine(t *testing.T) {
	issues := lintSample(t, "key=value\n\nkey2=value2\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoLeadingBlankLineReportsWithDuplicateRuleDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoDuplicateBlankLine = false

	issues := lintSampleWithConfig(t, "\nkey=value\n", config)

	assertIssueAt(t, issues, 1, 1, "Remove leading blank line.")
}

func TestNoLeadingBlankLineCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoLeadingBlankLine = false

	issues := lintSampleWithConfig(t, "\nkey=value\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
