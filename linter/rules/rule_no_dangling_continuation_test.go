package rules_test

import "testing"

func TestNoDanglingContinuationFlagsSingleTrailingBackslash(t *testing.T) {
	issues := lintSample(t, "foo=bar\\")

	assertIssueAt(t, issues, 1, 1, "Add line continuation or terminate.")
}

func TestNoDanglingContinuationAllowsEscapedBackslash(t *testing.T) {
	issues := lintSample(t, "foo=bar\\\\")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoDanglingContinuationAllowsContinuedEntry(t *testing.T) {
	issues := lintSample(t, "foo=bar\\\nbaz=qux\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoDanglingContinuationReportsEntryStartLine(t *testing.T) {
	issues := lintSample(t, "foo=bar\\\nbaz=qux\\")

	assertIssueAt(t, issues, 1, 1, "Add line continuation or terminate.")
}

func TestNoDanglingContinuationAllowsBlankContinuationTarget(t *testing.T) {
	issues := lintSample(t, "foo=bar\\\n\nbaz=qux\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoDanglingContinuationIgnoresComments(t *testing.T) {
	issues := lintSample(t, "# comment\\\nkey=value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoDanglingContinuationCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoDanglingContinuation = false

	issues := lintSampleWithConfig(t, "foo=bar\\", config)

	assertNoIssueContaining(t, issues, "Add line continuation or terminate.")
}
