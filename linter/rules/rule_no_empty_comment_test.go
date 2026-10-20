package rules_test

import "testing"

func TestNoEmptyCommentAllowsContent(t *testing.T) {
	issues := lintSample(t, "# comment\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoEmptyCommentAllowsSeparatorBetweenComments(t *testing.T) {
	issues := lintSample(t, "# one\n# \n# two\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoEmptyCommentFlagsStandaloneEmptyComment(t *testing.T) {
	issues := lintSample(t, "# \n")

	assertIssueAt(t, issues, 1, 1, "Remove empty '#'")
}

func TestNoEmptyCommentFlagsEmptyBlock(t *testing.T) {
	issues := lintSample(t, "#\n#\n")

	assertIssueAt(t, issues, 1, 1, "Remove empty '#'")
	assertIssueAt(t, issues, 2, 1, "Remove empty '#'")
}

func TestNoEmptyCommentIgnoresLeadingEmptyComment(t *testing.T) {
	issues := lintSample(t, "# \n# one\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
}

func TestNoEmptyCommentIgnoresTrailingEmptyComment(t *testing.T) {
	issues := lintSample(t, "# one\n# \n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
	assertIssueAt(t, issues, 2, 1, "Remove trailing empty comment.")
}

func TestNoEmptyCommentIgnoresEmptyCommentBeforeEntry(t *testing.T) {
	issues := lintSample(t, "# one\n# \nkey=value\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
	assertIssueAt(t, issues, 2, 1, "Remove trailing empty comment.")
}

func TestNoEmptyCommentIgnoresEmptyCommentBeforeBlankLine(t *testing.T) {
	issues := lintSample(t, "# one\n# \n\n# two\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
	assertIssueAt(t, issues, 2, 1, "Remove trailing empty comment.")
}

func TestNoEmptyCommentIgnoresConsecutiveEmptyComments(t *testing.T) {
	issues := lintSample(t, "# one\n# \n# \n# two\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoEmptyCommentFlagsEmptyBlockBeforeEntry(t *testing.T) {
	issues := lintSample(t, "#\nkey=value\n")

	assertIssueAt(t, issues, 1, 1, "Remove empty '#'")
}

func TestNoEmptyCommentFlagsEmptyBlockAfterContentBlock(t *testing.T) {
	issues := lintSample(t, "# one\n\n#\n")

	assertIssueAt(t, issues, 3, 1, "Remove empty '#'")
	assertNoIssueContaining(t, issues, "Remove trailing empty comment.")
}

func TestNoEmptyCommentAllowsMultipleSeparators(t *testing.T) {
	issues := lintSample(t, "# one\n# \n# two\n# \n# three\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestNoEmptyCommentAllowsWhitespaceOnlySeparator(t *testing.T) {
	issues := lintSample(t, "# one\n#    \n# two\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
}

func TestNoEmptyCommentFlagsEmptyBangComment(t *testing.T) {
	issues := lintSample(t, "! \n")

	assertIssueAt(t, issues, 1, 1, "Remove empty '#'")
}

func TestNoEmptyCommentAllowsBangSeparator(t *testing.T) {
	issues := lintSample(t, "! one\n! \n! two\n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
}

func TestNoEmptyCommentReportsCommentColumn(t *testing.T) {
	issues := lintSample(t, "  # \n#\n")

	assertIssueAt(t, issues, 1, 3, "Remove empty '#'")
	assertIssueAt(t, issues, 2, 1, "Remove empty '#'")
}

func TestNoEmptyCommentIgnoresInlineHash(t *testing.T) {
	issues := lintSample(t, "key=value # \n")

	assertNoIssueContaining(t, issues, "Remove empty '#'")
}

func TestNoEmptyCommentCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.CommentTrimming = false
	config.NoEmptyComment = false

	issues := lintSampleWithConfig(t, "# \n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
