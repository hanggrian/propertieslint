package propertieslint

func duplicateBlankLineIssue(path string, line int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Message: "duplicate blank line",
	}
}

func noLeadingBlankLineIssue(path string) *Issue {
	return &Issue{
		Path:    path,
		Line:    1,
		Message: "leading blank line",
	}
}

func trailingNewlineIssue(path string, lastLine int) *Issue {
	return &Issue{
		Path:    path,
		Line:    lastLine,
		Message: "missing trailing newline",
	}
}
