package rules_test

import "testing"

func TestValidEscapeAllowsKnownEscapes(t *testing.T) {
	issues := lintSample(t, `key=\t\n\r\f\\\ \:\=\#\!`+"\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestValidEscapeAllowsUnicodeEscapes(t *testing.T) {
	issues := lintSample(t, `key=\u0041\u00e9\u1234`+"\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestValidEscapeFlagsShortUnicodeInKey(t *testing.T) {
	issues := lintSample(t, "bad\\u12=value\n")

	assertIssueAt(t, issues, 1, 4, "Invalid key escape 'short unicode escape'.")
}

func TestValidEscapeFlagsShortUnicodeInValue(t *testing.T) {
	issues := lintSample(t, "value=\\u123\n")

	assertIssueAt(t, issues, 1, 7, "Invalid value escape 'short unicode escape'.")
}

func TestValidEscapeFlagsInvalidUnicodeDigits(t *testing.T) {
	issues := lintSample(t, "value=\\u12G4\n")

	assertIssueAt(t, issues, 1, 7, "Invalid value escape")
}

func TestValidEscapeFlagsTrailingBackslash(t *testing.T) {
	issues := lintSample(t, `key=value\\\`)

	assertIssueAt(t, issues, 1, 10, "Invalid value escape 'trailing backslash'.")
}

func TestValidEscapeInvalidKeyStopsLaterChecks(t *testing.T) {
	issues := lintSample(t, "bad\\u12=one\nbad\\u12=two\n")

	assertIssueAt(t, issues, 1, 4, "Invalid key escape")
	assertIssueAt(t, issues, 2, 4, "Invalid key escape")
	assertNoIssueContaining(t, issues, "duplicate key")
}

func TestValidEscapeAllowsEscapedSeparator(t *testing.T) {
	issues := lintSample(t, "key\\=part=value\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestValidEscapeCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.ValidEscape = false

	issues := lintSampleWithConfig(t, "bad\\u12=value\n", config)

	assertNoIssueContaining(t, issues, "Invalid key escape")
}
