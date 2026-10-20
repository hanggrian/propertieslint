package rules_test

import "testing"

func TestNoEmptyKeyFlagsEntryWithoutKey(t *testing.T) {
	issues := lintSample(t, "=bar\n")

	assertIssueAt(t, issues, 1, 1, "Add key.")
}

func TestNoEmptyKeyFlagsIndentedEntryWithoutKey(t *testing.T) {
	issues := lintSample(t, " =bar\n")

	assertIssueAt(t, issues, 1, 1, "Add key.")
}

func TestNoEmptyKeyReportsLogicalLine(t *testing.T) {
	issues := lintSample(t, "key=value\n=bar\n")

	assertIssueAt(t, issues, 2, 1, "Add key.")
}

func TestNoEmptyKeyReportsEachEmptyKey(t *testing.T) {
	issues := lintSample(t, "=1\n=2\n")

	assertIssueAt(t, issues, 1, 1, "Add key.")
	assertIssueAt(t, issues, 2, 1, "Add key.")
	assertNoIssueContaining(t, issues, "duplicate key")
}

func TestNoEmptyKeyIgnoresColonSeparator(t *testing.T) {
	issues := lintSample(t, ":value\n")

	assertNoIssueContaining(t, issues, "Add key.")
}

func TestNoEmptyKeyAllowsEmptyValue(t *testing.T) {
	issues := lintSample(t, "=\n")

	assertIssueAt(t, issues, 1, 1, "Add key.")
	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyKeyCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoEmptyKey = false

	issues := lintSampleWithConfig(t, "=bar\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
