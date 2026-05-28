package linter

import "strings"

func missingValueIssue(
	path string,
	line int,
	column int,
	separatorFound bool,
	value string,
) *Issue {
	if !separatorFound ||
		value != "" {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Missing value.",
	}
}

func quotedValueIssue(
	path string,
	line int,
	column int,
	value string,
) *Issue {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 2 ||
		trimmed[0] != '"' ||
		trimmed[len(trimmed)-1] != '"' {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "Value should not be quoted.",
	}
}
