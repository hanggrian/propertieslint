package rules_test

import "testing"

func TestCommentSpacingAllowsProperComment(t *testing.T) {
	issues := lintSample(t, "# comment\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentSpacingAllowsEmptyCommentWithSpace(t *testing.T) {
	issues := lintSample(t, "# one\n# \n# two\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentSpacingAllowsBareHash(t *testing.T) {
	issues := lintSample(t, "# one\n#\n# two\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentSpacingAllowsBareHashAtEndOfLine(t *testing.T) {
	config := DefaultConfig()
	config.CommentTrimming = false
	config.NoEmptyComment = false

	issues := lintSampleWithConfig(t, "#\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentSpacingFlagsMissingSpace(t *testing.T) {
	issues := lintSample(t, "#comment\n")

	assertIssueAt(t, issues, 1, 2, "Remove whitespace after '#'.")
}

func TestCommentSpacingFlagsExtraSpaces(t *testing.T) {
	issues := lintSample(t, "#  comment\n")

	assertIssueAt(t, issues, 1, 2, "Remove whitespace after '#'.")
}

func TestCommentSpacingFlagsLeadingWhitespace(t *testing.T) {
	issues := lintSample(t, " # comment\n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before '#'.")
}

func TestCommentSpacingFlagsBothSides(t *testing.T) {
	issues := lintSample(t, "  #comment\n")

	assertIssueAt(t, issues, 1, 1, "Remove whitespace before '#'.")
	assertIssueAt(t, issues, 1, 4, "Remove whitespace after '#'.")
}

func TestCommentSpacingFlagsDoubleHash(t *testing.T) {
	issues := lintSample(t, "## comment\n")

	assertIssueAt(t, issues, 1, 2, "Remove whitespace after '#'.")
}

func TestCommentSpacingIgnoresInlineHash(t *testing.T) {
	issues := lintSample(t, "key=value # comment\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentSpacingCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.CommentSpacing = false

	issues := lintSampleWithConfig(t, " #comment\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
