package linter

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Issue struct {
	Path    string
	Line    int
	Column  int
	Message string
}

type Result struct {
	CheckedFiles int
	Issues       []Issue
}

func Targets(paths []string, config Config) (Result, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	result := Result{}
	for _, target := range paths {
		if err := walkTarget(target, func(path string) error {
			issues, err := LintFile(path, config)
			if err != nil {
				return err
			}
			result.CheckedFiles++
			result.Issues = append(result.Issues, issues...)
			return nil
		}); err != nil {
			return Result{}, err
		}
	}

	sort.Slice(result.Issues, func(i, j int) bool {
		if result.Issues[i].Path != result.Issues[j].Path {
			return result.Issues[i].Path < result.Issues[j].Path
		}
		if result.Issues[i].Line != result.Issues[j].Line {
			return result.Issues[i].Line < result.Issues[j].Line
		}
		if result.Issues[i].Column != result.Issues[j].Column {
			return result.Issues[i].Column < result.Issues[j].Column
		}
		return result.Issues[i].Message < result.Issues[j].Message
	})

	return result, nil
}

func LintFile(path string, config Config) ([]Issue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	issues, err := lintReader(path, file, config)
	if err != nil {
		return nil, err
	}

	if config.TrailingNewline {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 &&
			data[len(data)-1] != '\n' {
			lineCount := strings.Count(string(data), "\n") + 1
			issues = append(issues, *trailingNewlineIssue(path, lineCount))
		}
	}

	return issues, nil
}

func lintReader(path string, r io.Reader, config Config) ([]Issue, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	seenKeys := map[string]int{}
	issues := make([]Issue, 0)
	logical := strings.Builder{}
	logicalRaw := strings.Builder{}
	logicalStartLine := 0
	continuing := false
	lineNumber := 0
	lastBlankLine := -2
	isFirstLine := true
	flush :=
		func(rawLogical string) error {
			// ensure flush only attempts key/value parsing for actual entries
			if logical.Len() == 0 {
				return nil
			}
			line := logical.String()
			logical.Reset()
			trimmed := strings.TrimLeft(line, " \t\f")
			if trimmed == "" {
				return nil
			}
			if trimmed[0] == '#' || trimmed[0] == '!' {
				return nil
			}

			// parse key and value
			key, value, separatorFound := splitKeyValue(trimmed)
			if key == "" {
				issues =
					append(issues, Issue{
						Path:    path,
						Line:    logicalStartLine,
						Column:  1,
						Message: "missing key",
					})
				return nil
			}

			// checks for violation
			if config.MissingSeparator {
				if issue :=
					missingSeparatorIssue(
						path,
						logicalStartLine,
						firstNonWhitespaceColumn(rawLogical),
						separatorFound,
					); issue != nil {
					issues = append(issues, *issue)
				}
			}
			normalizedKey := key
			if config.InvalidEscape {
				validatedKey, escapeIssues, valid :=
					validateEscapes(
						path,
						logicalStartLine,
						rawLogical,
						key,
						value,
					)
				issues = append(issues, escapeIssues...)
				if !valid {
					return nil
				}
				normalizedKey = validatedKey
			}
			if config.QuotedValue {
				if issue := quotedValueIssue(
					path,
					logicalStartLine,
					len(rawKeyStart(rawLogical))+2+firstQuoteColumn(rawValueStart(rawLogical))-1,
					value,
				); issue != nil {
					issues = append(issues, *issue)
				}
			}
			if config.KeyName {
				if issue := keyNameIssue(
					path,
					logicalStartLine,
					firstNonWhitespaceColumn(rawKeyStart(rawLogical)),
					normalizedKey,
				); issue != nil {
					issues = append(issues, *issue)
				}
			}
			if config.MissingValue {
				if issue :=
					missingValueIssue(
						path,
						logicalStartLine,
						len(rawKeyStart(rawLogical))+2,
						separatorFound,
						value,
					); issue != nil {
					issues = append(issues, *issue)
				}
			}
			if config.DuplicateKey {
				if issue :=
					duplicateKeyIssue(
						path,
						logicalStartLine,
						firstNonWhitespaceColumn(rawKeyStart(rawLogical)),
						normalizedKey,
						seenKeys,
					); issue != nil {
					issues = append(issues, *issue)
				}
			} else {
				seenKeys[normalizedKey] = logicalStartLine
			}
			if config.UntrimmedEntry {
				origKey, origValue, _ := splitKeyValueRaw(rawLogical)
				for _, issue := range untrimmedEntryIssue(
					path,
					logicalStartLine,
					origKey,
					origValue,
				) {
					issues = append(issues, *issue)
				}
			}

			// for inline comments, require one space around
			if config.CommentSpaces {
				line := rawLogical
				idx := -1
				for i := 0; i < len(line); i++ {
					if line[i] == '#' {
						// count preceding backslashes
						bs := 0
						j := i - 1
						for j >= 0 &&
							line[j] == '\\' {
							bs++
							j--
						}
						if bs%2 == 0 {
							idx = i
							break
						}
					}
				}
				if idx != -1 {
					beforeOK :=
						idx > 0 &&
							line[idx-1] == ' ' &&
							(idx-2 < 0 || line[idx-2] != ' ')
					afterOK :=
						idx+1 < len(line) &&
							line[idx+1] == ' ' &&
							(idx+2 >= len(line) || line[idx+2] != ' ')
					if !(beforeOK && afterOK) {
						issues =
							append(
								issues,
								*commentSpacesInlineIssue(path, logicalStartLine, idx+1),
							)
					}
				}
			}
			return nil
		}

	for scanner.Scan() {
		lineNumber++
		rawLine := scanner.Text()

		if !continuing {
			trimmed := strings.TrimLeft(rawLine, " \t\f")

			if trimmed == "" {
				if config.DuplicateBlankLine &&
					isFirstLine {
					if config.NoLeadingBlankLine {
						issues = append(issues, *noLeadingBlankLineIssue(path))
					}
				}
				if config.DuplicateBlankLine &&
					lastBlankLine == lineNumber-1 {
					issues = append(issues, *duplicateBlankLineIssue(path, lineNumber))
				}
				lastBlankLine = lineNumber
				isFirstLine = false
				continue
			}

			// for full-line comments, require zero left padding and exactly one space
			if trimmed[0] == '#' ||
				trimmed[0] == '!' {
				if trimmed[0] == '!' &&
					config.CommentStyle {
					issues =
						append(
							issues,
							*commentStyleIssue(path, lineNumber, firstNonWhitespaceColumn(rawLine)),
						)
				}
				if trimmed[0] == '#' &&
					config.CommentSpaces {
					leading := firstNonWhitespaceColumn(rawLine) - 1
					afterOK :=
						len(trimmed) >= 2 &&
							trimmed[1] == ' ' &&
							(len(trimmed) == 2 || trimmed[2] != ' ')
					if leading != 0 || !afterOK {
						issues =
							append(
								issues,
								*commentSpacesFullLineIssue(
									path,
									lineNumber,
									firstNonWhitespaceColumn(rawLine),
								),
							)
					}
				}
				isFirstLine = false
				continue
			}
			isFirstLine = false
			logicalStartLine = lineNumber
		} else {
			rawLine = strings.TrimLeft(rawLine, " \t\f")
		}

		// check for unterminated continuation from previous line
		if endsWithContinuation(rawLine) {
			logical.WriteString(rawLine[:len(rawLine)-1])
			logicalRaw.WriteString(rawLine[:len(rawLine)-1])
			continuing = true
			continue
		}

		// flush previous logical line if any
		logical.WriteString(rawLine)
		logicalRaw.WriteString(rawLine)
		continuing = false
		if err := flush(logicalRaw.String()); err != nil {
			return nil, err
		}
		logicalRaw.Reset()
	}

	// check for unterminated continuation at end of file
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if config.UnterminatedContinuation {
		if issue :=
			unterminatedContinuationIssue(
				path,
				logicalStartLine,
				continuing,
			); issue != nil {
			issues = append(issues, *issue)
		}
	}
	return issues, nil
}

func walkTarget(target string, visit func(path string) error) error {
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return visit(target)
	}
	return filepath.WalkDir(target, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(entry.Name()) != ".properties" {
			return nil
		}
		return visit(path)
	})
}

func isPropertiesWhitespace(ch byte) bool {
	return ch == ' ' ||
		ch == '\t' ||
		ch == '\f'
}

func firstWhitespaceColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if isPropertiesWhitespace(text[index]) {
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

func firstNonWhitespaceColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if !isPropertiesWhitespace(text[index]) {
			return index + 1
		}
	}
	return 1
}

func lastNonWhitespaceColumn(text string) int {
	for index := len(text) - 1; index >= 0; index-- {
		if !isPropertiesWhitespace(text[index]) {
			return index + 1
		}
	}
	return 1
}

func rawKeyStart(rawLogical string) string {
	key, _, _ := splitKeyValueRaw(rawLogical)
	return key
}

func rawValueStart(rawLogical string) string {
	_, value, _ := splitKeyValueRaw(rawLogical)
	return value
}

func firstQuoteColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if text[index] == '"' {
			return index + 1
		}
	}
	return 1
}
