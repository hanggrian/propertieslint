package propertieslint

import "fmt"

func validateEscapes(path string, line int, key string, value string) (string, []Issue, bool) {
	normalizedKey, err := unescape(key)
	if err != nil {
		return "", []Issue{{
			Path:    path,
			Line:    line,
			Message: fmt.Sprintf("invalid key escape: %v", err),
		}}, false
	}
	if _, err := unescape(value); err != nil {
		return normalizedKey, []Issue{{
			Path:    path,
			Line:    line,
			Message: fmt.Sprintf("invalid value escape: %v", err),
		}}, true
	}
	return normalizedKey, nil, true
}
