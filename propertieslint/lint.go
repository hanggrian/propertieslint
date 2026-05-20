package propertieslint

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
		if len(data) > 0 && data[len(data)-1] != '\n' {
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

	flush := func(rawLogical string) error {
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

		key, value, separatorFound := splitKeyValue(trimmed)
		if key == "" {
			issues = append(issues, Issue{
				Path:    path,
				Line:    logicalStartLine,
				Message: "missing key",
			})
			return nil
		}

		if config.MissingSeparator {
			if issue := missingSeparatorIssue(path, logicalStartLine, separatorFound); issue != nil {
				issues = append(issues, *issue)
			}
		}

		normalizedKey := key
		if config.InvalidEscape {
			validatedKey, escapeIssues, valid := validateEscapes(path, logicalStartLine, key, value)
			issues = append(issues, escapeIssues...)
			if !valid {
				return nil
			}
			normalizedKey = validatedKey
		}

		if config.MissingValue {
			if issue := missingValueIssue(path, logicalStartLine, separatorFound, value); issue != nil {
				issues = append(issues, *issue)
			}
		}

		if config.DuplicateKey {
			if issue := duplicateKeyIssue(path, logicalStartLine, normalizedKey, seenKeys); issue != nil {
				issues = append(issues, *issue)
			}
		} else {
			seenKeys[normalizedKey] = logicalStartLine
		}

		if config.UntrimmedEntry {
			origKey, origValue, _ := splitKeyValueRaw(rawLogical)
			for _, issue := range untrimmedEntryIssue(path, logicalStartLine, origKey, origValue) {
				issues = append(issues, *issue)
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
				if config.DuplicateBlankLine && isFirstLine {
					if config.NoLeadingBlankLine {
						issues = append(issues, *noLeadingBlankLineIssue(path))
					}
				}
				if config.DuplicateBlankLine && lastBlankLine == lineNumber-1 {
					issues = append(issues, *duplicateBlankLineIssue(path, lineNumber))
				}
				lastBlankLine = lineNumber
				isFirstLine = false
				continue
			}

			if trimmed[0] == '#' || trimmed[0] == '!' {
				isFirstLine = false
				continue
			}

			isFirstLine = false
			logicalStartLine = lineNumber
		} else {
			rawLine = strings.TrimLeft(rawLine, " \t\f")
		}

		if endsWithContinuation(rawLine) {
			logical.WriteString(rawLine[:len(rawLine)-1])
			logicalRaw.WriteString(rawLine[:len(rawLine)-1])
			continuing = true
			continue
		}

		logical.WriteString(rawLine)
		logicalRaw.WriteString(rawLine)
		continuing = false

		if err := flush(logicalRaw.String()); err != nil {
			return nil, err
		}
		logicalRaw.Reset()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if config.UnterminatedContinuation {
		if issue := unterminatedContinuationIssue(path, logicalStartLine, continuing); issue != nil {
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
