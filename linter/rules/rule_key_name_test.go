package rules_test

import "testing"

func TestKeyNameFlagsAllUppercase(t *testing.T) {
	issues := lintSample(t, "UPPER=value\n")

	assertIssueAt(t, issues, 1, 1, "Key name cannot be all uppercase.")
}

func TestKeyNameReportsKeyStartColumn(t *testing.T) {
	issues := lintSample(t, "  UPPER=value\n")

	assertIssueAt(t, issues, 1, 3, "Key name cannot be all uppercase.")
}

func TestKeyNameFlagsUppercaseWithDigits(t *testing.T) {
	issues := lintSample(t, "ITEM_2=value\n")

	assertIssueAt(t, issues, 1, 1, "Key name cannot be all uppercase.")
}

func TestKeyNameAllowsLowercaseNames(t *testing.T) {
	issues := lintSample(t, "item=value\nItem=value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyNameAllowsMixedAndSymbolHeavyNames(t *testing.T) {
	issues := lintSample(t, "api-V2=value\n123_456=value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyNameUsesUnescapedKey(t *testing.T) {
	issues := lintSample(t, "\\u0041BC=value\n")

	assertIssueAt(t, issues, 1, 1, "Key name cannot be all uppercase.")
}

func TestKeyNameAllowsEscapedLowercase(t *testing.T) {
	issues := lintSample(t, "A\\u0062C=value\n")

	assertNoIssueContaining(t, issues, "all uppercase")
}

func TestKeyNameWithoutEscapeNormalization(t *testing.T) {
	config := DefaultConfig()
	config.ValidEscape = false

	issues := lintSampleWithConfig(t, "\\u0041BC=value\n", config)

	assertNoIssueContaining(t, issues, "all uppercase")
}

func TestKeyNameHandlesUnicodeLetters(t *testing.T) {
	issues := lintSample(t, "ÄBC=value\näbc=value\n")

	assertIssueAt(t, issues, 1, 1, "Key name cannot be all uppercase.")
	assertNoIssueContaining(t, issues[1:], "all uppercase")
}

func TestKeyNameCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.KeyName = false

	issues := lintSampleWithConfig(t, "UPPER=value\n", config)

	assertNoIssueContaining(t, issues, "all uppercase")
}
