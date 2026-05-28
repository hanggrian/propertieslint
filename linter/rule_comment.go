package linter

func commentStyleIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Illegal comment '!'.",
	}
}

func commentSpacesFullLineIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "No space before and one space after '#'.",
	}
}

func commentSpacesInlineIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "One space around inline '#'.",
	}
}
