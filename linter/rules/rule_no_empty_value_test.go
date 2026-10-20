package rules_test

import "testing"

func noEmptyValueConfig() Config {
	config := DefaultConfig()
	config.NoEmptyValue = true
	return config
}

func TestNoEmptyValueFlagsEntryWithoutValue(t *testing.T) {
	issues := lintSampleWithConfig(t, "key=\n", noEmptyValueConfig())

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueReportsColumnAfterSeparator(t *testing.T) {
	issues := lintSampleWithConfig(t, "empty=\n", noEmptyValueConfig())

	assertIssueAt(t, issues, 1, 7, "Add value.")
}

func TestNoEmptyValueFlagsWhitespaceOnlyValue(t *testing.T) {
	issues := lintSampleWithConfig(t, "key=   \n", noEmptyValueConfig())

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueUsesRawKeyLength(t *testing.T) {
	issues := lintSampleWithConfig(t, "key = \n", noEmptyValueConfig())

	assertIssueAt(t, issues, 1, 6, "Add value.")
}

func TestNoEmptyValueAllowsEqualsAsValue(t *testing.T) {
	issues := lintSampleWithConfig(t, "key==\n", noEmptyValueConfig())

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueAllowsEscapedSpace(t *testing.T) {
	issues := lintSampleWithConfig(t, "key=\\ \n", noEmptyValueConfig())

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueIgnoresEmptyKey(t *testing.T) {
	issues := lintSampleWithConfig(t, "= \n", noEmptyValueConfig())

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueIgnoresMissingSeparator(t *testing.T) {
	issues := lintSampleWithConfig(t, "key\n", noEmptyValueConfig())

	assertNoIssueContaining(t, issues, "Add value")
}

func TestNoEmptyValueFlagsContinuedEmptyValue(t *testing.T) {
	issues := lintSampleWithConfig(t, "key=\\\n\n", noEmptyValueConfig())

	assertIssueAt(t, issues, 1, 5, "Add value.")
}

func TestNoEmptyValueCanBeDisabled(t *testing.T) {
	config := noEmptyValueConfig()
	config.NoEmptyValue = false

	issues := lintSampleWithConfig(t, "key=\n", config)

	assertNoIssueContaining(t, issues, "Add value")
}
