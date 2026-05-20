package propertieslint

func missingSeparatorIssue(path string, line int, separatorFound bool) *Issue {
	if separatorFound {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Message: "missing key/value separator",
	}
}
