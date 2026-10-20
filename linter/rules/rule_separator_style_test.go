package rules_test

import "testing"

func TestSeparatorStyleAllowsEquals(t *testing.T) {
	issues := lintSample(t, "key=value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestSeparatorStyleFlagsColon(t *testing.T) {
	issues := lintSample(t, "key:value\n")

	assertIssueAt(t, issues, 1, 1, "Use separator '='.")
}

func TestSeparatorStyleFlagsSpace(t *testing.T) {
	issues := lintSample(t, "key value\n")

	assertIssueAt(t, issues, 1, 1, "Use separator '='.")
}

func TestSeparatorStyleFlagsMissingSeparator(t *testing.T) {
	issues := lintSample(t, "keyOnly\n")

	assertIssueAt(t, issues, 1, 1, "Use separator '='.")
}

func TestSeparatorStyleReportsKeyStartColumn(t *testing.T) {
	issues := lintSample(t, "  key:value\n")

	assertIssueAt(t, issues, 1, 3, "Use separator '='.")
}

func TestSeparatorStyleAllowsEscapedEqualsInKey(t *testing.T) {
	issues := lintSample(t, "key\\=part=value\n")

	assertNoIssueContaining(t, issues, "Use separator '='.")
}

func TestSeparatorStyleUsesFirstEquals(t *testing.T) {
	issues := lintSample(t, "key==value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestSeparatorStyleDoesNotFlagEmptyKey(t *testing.T) {
	issues := lintSample(t, "=value\n")

	assertNoIssueContaining(t, issues, "Use separator '='.")
}

func TestSeparatorStyleFlagsTabSeparatedEntry(t *testing.T) {
	issues := lintSample(t, "key\tvalue\n")

	assertIssueAt(t, issues, 1, 1, "Use separator '='.")
}

func TestSeparatorStyleCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.SeparatorStyle = false

	issues := lintSampleWithConfig(t, "key:value\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
