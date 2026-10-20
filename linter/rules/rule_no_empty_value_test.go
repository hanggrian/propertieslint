package rules_test

import "testing"

func TestNoEmptyValueFlagsEntryWithoutValue(t *testing.T) {
	issues := lintSample(t, "key=\n")

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueReportsColumnAfterSeparator(t *testing.T) {
	issues := lintSample(t, "empty=\n")

	assertIssueAt(t, issues, 1, 7, "Add value.")
}

func TestNoEmptyValueFlagsWhitespaceOnlyValue(t *testing.T) {
	issues := lintSample(t, "key=   \n")

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueUsesRawKeyLength(t *testing.T) {
	issues := lintSample(t, "key = \n")

	assertIssueAt(t, issues, 1, 6, "Add value.")
}

func TestNoEmptyValueAllowsEqualsAsValue(t *testing.T) {
	issues := lintSample(t, "key==\n")

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueAllowsEscapedSpace(t *testing.T) {
	issues := lintSample(t, "key=\\ \n")

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueIgnoresEmptyKey(t *testing.T) {
	issues := lintSample(t, "= \n")

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueIgnoresMissingSeparator(t *testing.T) {
	issues := lintSample(t, "key\n")

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueFlagsContinuedEmptyValue(t *testing.T) {
	issues := lintSample(t, "key=\\\n\n")

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoEmptyValue = false

	issues := lintSampleWithConfig(t, "key=\n", config)

	assertNoIssueContaining(t, issues, "Add value")
}
