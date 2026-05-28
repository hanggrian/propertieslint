package linter

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultConfigFile = "propertieslint.json"

type Config struct {
	CommentStyle             bool
	CommentSpaces            bool
	DuplicateBlankLine       bool
	DuplicateKey             bool
	InvalidEscape            bool
	KeyName                  bool
	MissingSeparator         bool
	MissingValue             bool
	NoLeadingBlankLine       bool
	QuotedValue              bool
	TrailingNewline          bool
	UnterminatedContinuation bool
	UntrimmedEntry           bool
}

func DefaultConfig() Config {
	return Config{
		CommentStyle:             true,
		CommentSpaces:            true,
		DuplicateBlankLine:       true,
		DuplicateKey:             true,
		InvalidEscape:            true,
		KeyName:                  true,
		MissingSeparator:         true,
		MissingValue:             true,
		NoLeadingBlankLine:       true,
		QuotedValue:              true,
		TrailingNewline:          true,
		UnterminatedContinuation: true,
		UntrimmedEntry:           true,
	}
}

func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()
	if path == "" {
		return config, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	rules := map[string]bool{}
	if err := json.Unmarshal(data, &rules); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	for name, enabled := range rules {
		switch name {
		case "comment-style":
			config.CommentStyle = enabled
		case "comment-spaces":
			config.CommentSpaces = enabled
		case "duplicate-blank-line":
			config.DuplicateBlankLine = enabled
		case "duplicate-key":
			config.DuplicateKey = enabled
		case "invalid-escape":
			config.InvalidEscape = enabled
		case "key-name":
			config.KeyName = enabled
		case "missing-key-value-separator":
			config.MissingSeparator = enabled
		case "missing-value":
			config.MissingValue = enabled
		case "no-leading-blank-line":
			config.NoLeadingBlankLine = enabled
		case "quoted-value":
			config.QuotedValue = enabled
		case "trailing-newline":
			config.TrailingNewline = enabled
		case "unterminated-line-continuation":
			config.UnterminatedContinuation = enabled
		case "untrimmed-entry":
			config.UntrimmedEntry = enabled
		default:
			return Config{}, fmt.Errorf("Unknown rule %q", name)
		}
	}
	return config, nil
}

func ResolveConfigPath(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		return DefaultConfigFile
	}
	return ""
}
