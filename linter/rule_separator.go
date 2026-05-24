package linter

func missingSeparatorIssue(
	path string,
	line int,
	column int,
	separatorFound bool,
) *Issue {
	if separatorFound {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Column:  column,
		Message: "missing key/value separator",
	}
}
