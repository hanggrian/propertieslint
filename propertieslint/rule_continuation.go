package propertieslint

func unterminatedContinuationIssue(path string, line int, continuing bool) *Issue {
	if !continuing {
		return nil
	}
	return &Issue{
		Path:    path,
		Line:    line,
		Message: "unterminated line continuation",
	}
}
