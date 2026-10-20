package linter

import (
	"strings"

	"github.com/hanggrian/propertieslint/linter/rules"
)

func splitKeyValue(line string) (key string, value string, separatorFound bool) {
	escaped := false
	for index := 0; index < len(line); index++ {
		current := line[index]
		if escaped {
			escaped = false
			continue
		}
		if current == '\\' {
			escaped = true
			continue
		}
		if current == '=' {
			return trimUnescapedRight(line[:index]),
				strings.TrimLeft(line[index+1:], " \t\f"),
				true
		}
	}
	return trimUnescapedRight(strings.TrimLeft(line, " \t\f")), "", false
}

func trimUnescapedRight(text string) string {
	end := len(text)
	for end > 0 && rules.IsPropertiesWhitespace(text[end-1]) {
		backslashes := 0
		for index := end - 2; index >= 0 && text[index] == '\\'; index-- {
			backslashes++
		}
		if backslashes%2 == 1 {
			break
		}
		end--
	}
	return text[:end]
}

func endsWithContinuation(line string) bool {
	if len(line) == 0 {
		return false
	}
	backslashes := 0
	for index := len(line) - 1; index >= 0; index-- {
		if line[index] != '\\' {
			break
		}
		backslashes++
	}
	return backslashes == 1
}
