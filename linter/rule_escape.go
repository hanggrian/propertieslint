package linter

import "fmt"

func validateEscapes(
	path string,
	line int,
	rawLogical string,
	key string,
	value string,
) (string, []Issue, bool) {
	rawKey, _, _ := splitKeyValueRaw(rawLogical)
	normalizedKey, err := unescape(key)
	if err != nil {
		return "", []Issue{{
			Path:    path,
			Line:    line,
			Column:  firstBackslashColumn(rawKey),
			Message: fmt.Sprintf("invalid key escape: %v", err),
		}}, false
	}
	if _, err := unescape(value); err != nil {
		return normalizedKey, []Issue{{
			Path:    path,
			Line:    line,
			Column:  len(rawKey) + 2 + firstBackslashColumn(value) - 1,
			Message: fmt.Sprintf("invalid value escape: %v", err),
		}}, true
	}
	return normalizedKey, nil, true
}
