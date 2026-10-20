package rules

func CommentStyleIssue(path string, line int, column int) *Issue {
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Use comment '#'.",
	}
}

func CommentSpacingIssue(
	path string,
	line int,
	column int,
	trimmed string,
) []*Issue {
	var issues []*Issue
	if column > 1 {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  1,
			Message: "Remove whitespace before '#'.",
		})
	}
	afterOK :=
		len(trimmed) >= 2 &&
			trimmed[1] == ' ' &&
			(len(trimmed) == 2 || trimmed[2] != ' ')
	if !afterOK {
		issues = append(issues, &Issue{
			Path:    path,
			Line:    line,
			Column:  column + 1,
			Message: "Remove whitespace after '#'.",
		})
	}
	return issues
}
