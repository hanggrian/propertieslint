package rules_test

import "testing"

func commentTrimmingConfig() Config {
	config := DefaultConfig()
	config.NoEmptyComment = false
	return config
}

func TestCommentTrimmingAllowsContentOnly(t *testing.T) {
	issues := lintSampleWithConfig(t, "# one\n# two\n", commentTrimmingConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingAllowsInteriorEmptyComment(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# one\n#\n# two\n",
		commentTrimmingConfig(),
	)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingFlagsLeadingEmptyComment(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"#\n# one\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
}

func TestCommentTrimmingFlagsTrailingEmptyComment(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# one\n#\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 2, 1, "Remove trailing empty comment.")
}

func TestCommentTrimmingFlagsBothEdges(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"#\n# one\n#\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
	assertIssueAt(t, issues, 3, 1, "Remove trailing empty comment.")
}

func TestCommentTrimmingFlagsMultipleLeadingComments(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"#\n#\n# one\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
	assertIssueAt(t, issues, 2, 1, "Remove leading empty comment.")
}

func TestCommentTrimmingFlagsMultipleTrailingComments(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# one\n#\n#\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 2, 1, "Remove trailing empty comment.")
	assertIssueAt(t, issues, 3, 1, "Remove trailing empty comment.")
}

func TestCommentTrimmingIgnoresEmptyBlock(t *testing.T) {
	issues := lintSampleWithConfig(t, "#\n#\n", commentTrimmingConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingAllowsStandaloneEmptyComment(t *testing.T) {
	issues := lintSampleWithConfig(t, "#\n", commentTrimmingConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingAllowsSeparateEmptyBlock(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# one\n\n#\n",
		commentTrimmingConfig(),
	)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingAllowsLeadingEmptyBlock(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"#\n\n# one\n",
		commentTrimmingConfig(),
	)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingAllowsEmptyBlockAfterEntry(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# one\nkey=value\n#\n",
		commentTrimmingConfig(),
	)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestCommentTrimmingReportsCommentColumn(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"  #\n# one\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 3, "Remove leading empty comment.")
}

func TestCommentTrimmingFlagsBangComments(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"!\n! one\n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
}

func TestCommentTrimmingAllowsInteriorButFlagsEdges(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"#\n# one\n#\n# two\n#\n",
		commentTrimmingConfig(),
	)

	if len(issues) != 2 {
		t.Fatalf("expected two issues, got %#v", issues)
	}
	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
	assertIssueAt(t, issues, 5, 1, "Remove trailing empty comment.")
}

func TestCommentTrimmingFlagsWhitespaceOnlyComments(t *testing.T) {
	issues := lintSampleWithConfig(
		t,
		"# \n# one\n# \n",
		commentTrimmingConfig(),
	)

	assertIssueAt(t, issues, 1, 1, "Remove leading empty comment.")
	assertIssueAt(t, issues, 3, 1, "Remove trailing empty comment.")
}

func TestCommentTrimmingCanBeDisabled(t *testing.T) {
	config := commentTrimmingConfig()
	config.CommentTrimming = false

	issues := lintSampleWithConfig(t, "#\n# one\n#\n", config)

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}
