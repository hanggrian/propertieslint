package rules_test

import "testing"

func TestCommentStyleFlagsBangComment(t *testing.T) {
	issues := lintSample(t, "!notacomment\n")

	assertIssueAt(t, issues, 1, 1, "Use comment '#'.")
}

func TestCommentStyleFlagsIndentedBangComment(t *testing.T) {
	issues := lintSample(t, "  !notacomment\n")

	assertIssueAt(t, issues, 1, 3, "Use comment '#'.")
}

func TestCommentStyleAllowsHashComment(t *testing.T) {
	issues := lintSample(t, "# comment\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentStyleAllowsBangInsideValue(t *testing.T) {
	issues := lintSample(t, "key=value!more\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentStyleAllowsBangOnContinuationLine(t *testing.T) {
	issues := lintSample(t, "key=line1\\\n!line2\n")

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentStyleCanBeDisabled(t *testing.T) {
	config := DefaultConfig()
	config.CommentStyle = false

	issues := lintSampleWithConfig(t, "!notacomment\n", config)

	assertNoIssueContaining(t, issues, "Use comment '#'.")
}
