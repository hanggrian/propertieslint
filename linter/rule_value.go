package linter

func missingValueIssue(path string, line int, column int, separatorFound bool, value string) *Issue {
	if !separatorFound ||
		value != "" {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "missing value",
	}
}
