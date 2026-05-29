package linter

func duplicateBlankLineIssue(path string, line int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  1,
		Message: "Duplicate blank line.",
	}
}

func noLeadingBlankLineIssue(path string) *Issue {
	return &Issue{
		Path:    path,
		Line:    1,
		Column:  1,
		Message: "Unexpected leading blank line.",
	}
}

func trailingNewlineIssue(path string, lastLine int) *Issue {
	return &Issue{
		Path:    path,
		Line:    lastLine,
		Column:  1,
		Message: "Missing trailing newline",
	}
}
