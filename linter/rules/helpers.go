package rules

import (
	"errors"
	"fmt"
	"strings"
)

func IsPropertiesWhitespace(ch byte) bool {
	return ch == ' ' ||
		ch == '\t' ||
		ch == '\f'
}

func firstWhitespaceColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if IsPropertiesWhitespace(text[index]) {
			return index + 1
		}
	}
	return 1
}

func firstBackslashColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if text[index] == '\\' {
			return index + 1
		}
	}
	return 1
}

func lastNonWhitespaceColumn(text string) int {
	for index := len(text) - 1; index >= 0; index-- {
		if !IsPropertiesWhitespace(text[index]) {
			return index + 1
		}
	}
	return 1
}

func SplitKeyValueRaw(line string) (key string, value string, separatorFound bool) {
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
			return line[:index], line[index+1:], true
		}
	}
	return line, "", false
}

func unescape(value string) (string, error) {
	if value == "" {
		return value, nil
	}
	var builder strings.Builder
	for index := 0; index < len(value); index++ {
		current := value[index]
		if current != '\\' {
			builder.WriteByte(current)
			continue
		}

		index++
		if index >= len(value) {
			return "", errors.New("trailing backslash")
		}

		switch value[index] {
		case 't':
			builder.WriteByte('\t')

		case 'n':
			builder.WriteByte('\n')

		case 'r':
			builder.WriteByte('\r')

		case 'f':
			builder.WriteByte('\f')

		case '\\', ' ', ':', '=', '#', '!':
			builder.WriteByte(value[index])

		case 'u':
			if index+4 >= len(value) {
				return "", errors.New("short unicode escape")
			}
			codePoint := 0
			for offset := 1; offset <= 4; offset++ {
				hexDigit := value[index+offset]
				codePoint *= 16
				switch {
				case hexDigit >= '0' && hexDigit <= '9':
					codePoint += int(hexDigit - '0')
				case hexDigit >= 'a' && hexDigit <= 'f':
					codePoint += int(hexDigit-'a') + 10
				case hexDigit >= 'A' && hexDigit <= 'F':
					codePoint += int(hexDigit-'A') + 10
				default:
					return "", fmt.Errorf("invalid unicode escape %q", value[index:index+5])
				}
			}
			builder.WriteRune(rune(codePoint))
			index += 4

		default:
			builder.WriteByte(value[index])
		}
	}
	return builder.String(), nil
}
