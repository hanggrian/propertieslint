package rules_test

import "testing"

func TestNoDuplicateKeyFlagsConsecutiveDuplicate(t *testing.T) {
	issues := lintSample(t, "duplicate=one\nduplicate=two\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyReportsKeyStartColumn(t *testing.T) {
	issues := lintSample(t, "duplicate=one\n duplicate=two\n")

	assertIssueAt(t, issues, 2, 2, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyIsCaseSensitive(t *testing.T) {
	issues := lintSample(t, "item=1\nItem=2\n")

	assertNoIssueContaining(t, issues, "duplicate key")
}

func TestNoDuplicateKeyDetectsAfterCommentAndBlankLine(t *testing.T) {
	issues := lintSample(t, "item=1\n# comment\n\nitem=2\n")

	assertIssueAt(t, issues, 4, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyComparesUnescapedKeys(t *testing.T) {
	issues := lintSample(t, "a\\u0041=1\naA=2\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyComparesEscapedSeparatorKeys(t *testing.T) {
	issues := lintSample(t, "a\\=b=1\na\\=b=2\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyDetectsAcrossContinuation(t *testing.T) {
	issues := lintSample(t, "item=one\nitem\\\n=two\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyReportsEachDuplicateAgainstFirst(t *testing.T) {
	issues := lintSample(t, "dup=1\ndup=2\ndup=3\n")

	assertIssueAt(t, issues, 2, 1, "Remove duplicate key (first seen at line 1).")
	assertIssueAt(t, issues, 3, 1, "Remove duplicate key (first seen at line 1).")
}

func TestNoDuplicateKeyIgnoresEmptyKeys(t *testing.T) {
	issues := lintSample(t, "=1\n=2\n")

	assertNoIssueContaining(t, issues, "duplicate key")
}

func TestNoDuplicateKeyWithoutEscapeNormalization(t *testing.T) {
	config := DefaultConfig()
	config.ValidEscape = false

	issues := lintSampleWithConfig(t, "a\\u0041=1\naA=2\n", config)

	assertNoIssueContaining(t, issues, "duplicate key")
}

func TestNoDuplicateKeyCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.NoDuplicateKey = false

	issues := lintSampleWithConfig(t, "duplicate=one\nduplicate=two\n", config)

	assertNoIssueContaining(t, issues, "duplicate key")
}
