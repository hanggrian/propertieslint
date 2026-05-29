package linter

import (
	"bytes"
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
	MissingKey               bool
	MissingSeparator         bool
	MissingValue             bool
	NoLeadingBlankLine       bool
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
		MissingKey:               true,
		MissingSeparator:         true,
		MissingValue:             true,
		NoLeadingBlankLine:       true,
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
	groups := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &groups); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	for groupName, rawGroup := range groups {
		if groupName != "comment" &&
			groupName != "format" &&
			groupName != "pair" &&
			groupName != "whitespace" {
			return Config{}, fmt.Errorf("unknown group %q in config %q", groupName, path)
		}
		if bytes.Equal(bytes.TrimSpace(rawGroup), []byte("false")) {
			disableGroup(&config, groupName)
			continue
		}
		if bytes.Equal(bytes.TrimSpace(rawGroup), []byte("true")) {
			return Config{}, fmt.Errorf(
				"parse config %q: group %q must be false or an object of rules",
				path,
				groupName,
			)
		}
		rules := map[string]bool{}
		if err := json.Unmarshal(rawGroup, &rules); err != nil {
			return Config{}, err
		}
		for ruleName, enabled := range rules {
			if err := applyGroupRule(&config, groupName, ruleName, enabled); err != nil {
				return Config{}, err
			}
		}
	}
	return config, nil
}

func applyGroupRule(config *Config, groupName string, ruleName string, enabled bool) error {
	switch groupName {
	case "comment":
		switch ruleName {
		case "comment-style":
			config.CommentStyle = enabled
		case "comment-spaces":
			config.CommentSpaces = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "format":
		switch ruleName {
		case "invalid-escape":
			config.InvalidEscape = enabled
		case "missing-key-value-separator":
			config.MissingSeparator = enabled
		case "unterminated-line-continuation":
			config.UnterminatedContinuation = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "pair":
		switch ruleName {
		case "duplicate-key":
			config.DuplicateKey = enabled
		case "key-name":
			config.KeyName = enabled
		case "missing-key":
			config.MissingKey = enabled
		case "missing-value":
			config.MissingValue = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "whitespace":
		switch ruleName {
		case "duplicate-blank-line":
			config.DuplicateBlankLine = enabled
		case "no-leading-blank-line":
			config.NoLeadingBlankLine = enabled
		case "trailing-newline":
			config.TrailingNewline = enabled
		case "untrimmed-entry":
			config.UntrimmedEntry = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	default:
		return fmt.Errorf("unknown group %q", groupName)
	}
	return nil
}

func disableGroup(config *Config, groupName string) {
	switch groupName {
	case "comment":
		config.CommentStyle = false
		config.CommentSpaces = false

	case "format":
		config.InvalidEscape = false
		config.MissingSeparator = false
		config.UnterminatedContinuation = false

	case "pair":
		config.DuplicateKey = false
		config.KeyName = false
		config.MissingKey = false
		config.MissingValue = false

	case "whitespace":
		config.DuplicateBlankLine = false
		config.NoLeadingBlankLine = false
		config.TrailingNewline = false
		config.UntrimmedEntry = false
	}
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
