package linter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const DefaultConfigFile = "propertieslint.json"

type Config struct {
	CommentStyle           bool
	CommentSpacing         bool
	EntryTrimming          bool
	KeyName                bool
	NoDanglingContinuation bool
	NoDuplicateBlankLine   bool
	NoDuplicateKey         bool
	NoEmptyKey             bool
	NoEmptyValue           bool
	NoLeadingBlankLine     bool
	SeparatorStyle         bool
	TrailingNewline        bool
	ValidEscape            bool
}

func DefaultConfig() Config {
	return Config{
		CommentStyle:           true,
		CommentSpacing:         true,
		EntryTrimming:          true,
		KeyName:                true,
		NoDanglingContinuation: true,
		NoDuplicateBlankLine:   true,
		NoDuplicateKey:         true,
		NoEmptyKey:             true,
		NoEmptyValue:           true,
		NoLeadingBlankLine:     true,
		SeparatorStyle:         true,
		TrailingNewline:        true,
		ValidEscape:            true,
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
		case "comment-spacing":
			config.CommentSpacing = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "format":
		switch ruleName {
		case "no-dangling-continuation":
			config.NoDanglingContinuation = enabled
		case "separator-style":
			config.SeparatorStyle = enabled
		case "valid-escape":
			config.ValidEscape = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "pair":
		switch ruleName {
		case "key-name":
			config.KeyName = enabled
		case "no-duplicate-key":
			config.NoDuplicateKey = enabled
		case "no-empty-key":
			config.NoEmptyKey = enabled
		case "no-empty-value":
			config.NoEmptyValue = enabled
		default:
			return fmt.Errorf("unknown rule %q in group %q", ruleName, groupName)
		}

	case "whitespace":
		switch ruleName {
		case "entry-trimming":
			config.EntryTrimming = enabled
		case "no-duplicate-blank-line":
			config.NoDuplicateBlankLine = enabled
		case "no-leading-blank-line":
			config.NoLeadingBlankLine = enabled
		case "trailing-newline":
			config.TrailingNewline = enabled
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
		config.CommentSpacing = false

	case "format":
		config.NoDanglingContinuation = false
		config.SeparatorStyle = false
		config.ValidEscape = false

	case "pair":
		config.KeyName = false
		config.NoDuplicateKey = false
		config.NoEmptyKey = false
		config.NoEmptyValue = false

	case "whitespace":
		config.EntryTrimming = false
		config.NoDuplicateBlankLine = false
		config.NoLeadingBlankLine = false
		config.TrailingNewline = false
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
