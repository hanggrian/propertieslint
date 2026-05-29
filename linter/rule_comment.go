package linter

func commentStyleIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Illegal comment '!'.",
	}
}

func commentSpacesIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "No space before and one space after '#'.",
	}
}
