package linter

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hanggrian/propertieslint/linter/rules"
)

type Issue = rules.Issue

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
		if err :=
			walkTarget(
				target,
				func(path string) error {
					issues, err := LintFile(path, config)
					if err != nil {
						return err
					}
					result.CheckedFiles++
					result.Issues = append(result.Issues, issues...)
					return nil
				},
			); err != nil {
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

	issues, err := LintReader(path, file, config)
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
			issues = append(issues, *rules.TrailingNewlineIssue(path, lineCount))
		}
	}

	return issues, nil
}

// LintReader reports rule violations in a properties stream.
func LintReader(path string, r io.Reader, config Config) ([]Issue, error) {
	type commentPosition struct {
		line   int
		column int
	}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	seenKeys := map[string]int{}
	previousKey := ""
	issues := make([]Issue, 0)
	logical := strings.Builder{}
	logicalRaw := strings.Builder{}
	logicalStartLine := 0
	continuing := false
	lineNumber := 0
	lastBlankLine := -2
	isFirstLine := true
	pendingEmptyComments := make([]commentPosition, 0)
	commentBlockHasContent := false
	flushTrimmedComments := func(leading bool) {
		if config.CommentTrimming {
			for _, comment := range pendingEmptyComments {
				issues = append(issues, *rules.CommentTrimmingIssue(
					path,
					comment.line,
					comment.column,
					leading,
				))
			}
		}
		pendingEmptyComments = pendingEmptyComments[:0]
	}
	endCommentBlock := func() {
		if commentBlockHasContent {
			flushTrimmedComments(false)
		} else if config.NoEmptyComment {
			for _, comment := range pendingEmptyComments {
				issues = append(issues, *rules.NoEmptyCommentIssue(
					path,
					comment.line,
					comment.column,
				))
			}
		}
		pendingEmptyComments = pendingEmptyComments[:0]
		commentBlockHasContent = false
	}
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
			if key == "" && separatorFound {
				if config.NoEmptyKey {
					issues = append(issues, *rules.NoEmptyKeyIssue(path, logicalStartLine))
				}
				return nil
			}

			// checks for violation
			if config.SeparatorStyle {
				if issue :=
					rules.SeparatorStyleIssue(
						path,
						logicalStartLine,
						firstNonWhitespaceColumn(rawLogical),
						separatorFound,
					); issue != nil {
					issues = append(issues, *issue)
				}
			}
			normalizedKey := key
			if config.ValidEscape {
				validatedKey, escapeIssues, valid :=
					rules.ValidateEscape(
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
			if config.KeyOrder {
				if issue := rules.KeyOrderIssue(
					path,
					logicalStartLine,
					firstNonWhitespaceColumn(rawKeyStart(rawLogical)),
					previousKey,
					normalizedKey,
				); issue != nil {
					issues = append(issues, *issue)
				}
			}
			previousKey = normalizedKey
			if config.KeyName {
				if issue := rules.KeyNameIssue(
					path,
					logicalStartLine,
					firstNonWhitespaceColumn(rawKeyStart(rawLogical)),
					normalizedKey,
				); issue != nil {
					issues = append(issues, *issue)
				}
			}
			if config.NoEmptyValue {
				if issue :=
					rules.NoEmptyValueIssue(
						path,
						logicalStartLine,
						len(rawKeyStart(rawLogical))+2,
						separatorFound,
						value,
					); issue != nil {
					issues = append(issues, *issue)
				}
			}
			if config.NoDuplicateKey {
				if issue :=
					rules.NoDuplicateKeyIssue(
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
			if config.EntryTrimming {
				origKey, origValue, _ := rules.SplitKeyValueRaw(rawLogical)
				for _, issue := range rules.EntryTrimmingIssue(
					path,
					logicalStartLine,
					origKey,
					origValue,
				) {
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
				if isFirstLine &&
					config.NoLeadingBlankLine {
					issues = append(issues, *rules.NoLeadingBlankLineIssue(path))
				}
				if config.NoDuplicateBlankLine &&
					lastBlankLine == lineNumber-1 {
					issues = append(issues, *rules.NoDuplicateBlankLineIssue(path, lineNumber))
				}
				lastBlankLine = lineNumber
				isFirstLine = false
				endCommentBlock()
				continue
			}

			// comment checks
			if trimmed[0] == '#' ||
				trimmed[0] == '!' {
				hasContent := strings.TrimSpace(trimmed[1:]) != ""
				if hasContent {
					if !commentBlockHasContent {
						flushTrimmedComments(true)
					} else {
						pendingEmptyComments = pendingEmptyComments[:0]
					}
					commentBlockHasContent = true
				} else {
					pendingEmptyComments = append(
						pendingEmptyComments,
						commentPosition{
							line:   lineNumber,
							column: firstNonWhitespaceColumn(rawLine),
						},
					)
				}
				if trimmed[0] == '!' &&
					config.CommentStyle {
					issues =
						append(
							issues,
							*rules.CommentStyleIssue(
								path,
								lineNumber,
								firstNonWhitespaceColumn(rawLine),
							),
						)
				}
				if trimmed[0] == '#' &&
					config.CommentSpacing {
					for _, issue := range rules.CommentSpacingIssue(
						path,
						lineNumber,
						firstNonWhitespaceColumn(rawLine),
						trimmed,
					) {
						issues = append(issues, *issue)
					}
				}
				isFirstLine = false
				continue
			}
			isFirstLine = false
			endCommentBlock()
			logicalStartLine = lineNumber
		} else {
			rawLine = strings.TrimLeft(rawLine, " \t\f")
		}

		// check whether this line continues onto the next
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

	// check for dangling continuation at end of file
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if config.NoDanglingContinuation {
		if issue :=
			rules.NoDanglingContinuationIssue(
				path,
				logicalStartLine,
				continuing,
			); issue != nil {
			issues = append(issues, *issue)
		}
	}
	endCommentBlock()
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
	return filepath.WalkDir(
		target,
		func(path string, entry os.DirEntry, err error) error {
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
		},
	)
}

func firstNonWhitespaceColumn(text string) int {
	for index := 0; index < len(text); index++ {
		if !rules.IsPropertiesWhitespace(text[index]) {
			return index + 1
		}
	}
	return 1
}

func rawKeyStart(rawLogical string) string {
	key, _, _ := rules.SplitKeyValueRaw(rawLogical)
	return key
}
