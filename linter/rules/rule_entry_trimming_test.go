package rules_test

import "testing"

func TestEntryTrimmingFlagsKeyLeadingWhitespace(t *testing.T) {
	issues := lintSample(t, " key=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before key.")
}

func TestEntryTrimmingFlagsKeyTrailingWhitespace(t *testing.T) {
	issues := lintSample(t, "key =value\n")

	assertIssueAt(t, issues, 1, 4, "Remove whitespace after key.")
}

func TestEntryTrimmingFlagsValueLeadingWhitespace(t *testing.T) {
	issues := lintSample(t, "key= value\n")

	assertIssueAt(t, issues, 1, 5, "Remove whitespace before value.")
}

func TestEntryTrimmingFlagsValueTrailingWhitespace(t *testing.T) {
	issues := lintSample(t, "key=value \n")

	assertIssueAt(t, issues, 1, 10, "Remove whitespace after value.")
}

func TestEntryTrimmingFlagsTabs(t *testing.T) {
	issues := lintSample(t, "\tkey=value\nkey\t=value\nkey=\tvalue\nkey=value\t\n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before key.")
	assertIssueAt(t, issues, 2, 4, "Remove whitespace after key.")
	assertIssueAt(t, issues, 3, 5, "Remove whitespace before value.")
	assertIssueAt(t, issues, 4, 10, "Remove whitespace after value.")
}

func TestEntryTrimmingFlagsFormFeed(t *testing.T) {
	issues := lintSample(t, "\fkey=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before key.")
}

func TestEntryTrimmingFlagsBothKeyAndValue(t *testing.T) {
	issues := lintSample(t, " key = value \n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before key.")
	assertIssueAt(t, issues, 1, 5, "Remove whitespace after key.")
	assertIssueAt(t, issues, 1, 7, "Remove whitespace before value.")
	assertIssueAt(t, issues, 1, 13, "Remove whitespace after value.")
}

func TestEntryTrimmingAllowsContinuationIndentation(t *testing.T) {
	issues := lintSample(t, "key=line1\\\n  line2\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestEntryTrimmingFlagsEscapedTrailingWhitespace(t *testing.T) {
	issues := lintSample(t, "key\\ =value\n")

	assertIssueAt(t, issues, 1, 5, "Remove whitespace after key.")
}

func TestEntryTrimmingIgnoresEmptyKeyAndValue(t *testing.T) {
	issues := lintSample(t, "=value\nkey=\n")

	assertNoIssueContaining(t, issues, "whitespace")
}

func TestEntryTrimmingCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.EntryTrimming = false

	issues := lintSampleWithConfig(t, " key = value \n", config)

	assertNoIssueContaining(t, issues, "whitespace")
}
