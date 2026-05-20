package propertieslint

func missingValueIssue(path string, line int, separatorFound bool, value string) *Issue {
	if !separatorFound || value != "" {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Message: "missing value",
	}
}
