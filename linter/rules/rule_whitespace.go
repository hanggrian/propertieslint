package rules

func NoDuplicateBlankLineIssue(path string, line int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  1,
		Message: "Remove duplicate blank line.",
	}
}

func NoLeadingBlankLineIssue(path string) *Issue {
	return &Issue{
		Path:    path,
		Line:    1,
		Column:  1,
		Message: "Remove leading blank line.",
	}
}

func TrailingNewlineIssue(path string, lastLine int) *Issue {
	return &Issue{
		Path:    path,
		Line:    lastLine,
		Column:  1,
		Message: "Add trailing newline",
	}
}
